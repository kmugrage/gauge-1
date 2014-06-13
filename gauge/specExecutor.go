package main

import (
	"code.google.com/p/goprotobuf/proto"
	"net"
)

type itemExecutor func(item, *specExecutor) *stepExecutionStatus

type specExecutor struct {
	specification        *specification
	dataTableIndex       int
	connection           net.Conn
	conceptDictionary    *conceptDictionary
	pluginHandler        *pluginHandler
	currentExecutionInfo *ExecutionInfo
	specResult           *specResult
}

func (specExecutor *specExecutor) initialize(specificationToExecute *specification, connection net.Conn, pluginHandler *pluginHandler) {
	specExecutor.specification = specificationToExecute
	specExecutor.connection = connection
	specExecutor.pluginHandler = pluginHandler
}

type stepValidationError struct {
	step     *step
	message  string
	fileName string
}

func (e *stepValidationError) Error() string {
	return e.message
}

func (e *specExecutor) executeBeforeSpecHook() *ProtoExecutionResult {
	message := &Message{MessageType: Message_SpecExecutionStarting.Enum(),
		SpecExecutionStartingRequest: &SpecExecutionStartingRequest{CurrentExecutionInfo: e.currentExecutionInfo}}

	e.pluginHandler.notifyPlugins(message)
	return executeAndGetStatus(e.connection, message)
}

func (e *specExecutor) executeAfterSpecHook() *ProtoExecutionResult {
	message := &Message{MessageType: Message_SpecExecutionEnding.Enum(),
		SpecExecutionEndingRequest: &SpecExecutionEndingRequest{CurrentExecutionInfo: e.currentExecutionInfo}}
	e.pluginHandler.notifyPlugins(message)
	return executeAndGetStatus(e.connection, message)
}

func (specExecutor *specExecutor) execute() *specResult {
	specInfo := &SpecInfo{Name: proto.String(specExecutor.specification.heading.value),
		FileName: proto.String(specExecutor.specification.fileName),
		IsFailed: proto.Bool(false), Tags: getTagValue(specExecutor.specification.tags)}
	specExecutor.currentExecutionInfo = &ExecutionInfo{CurrentSpec: specInfo}

	getCurrentConsole().writeSpecHeading(specExecutor.specification)

	specExecutor.specResult = &specResult{}
	specExecutor.specResult.addSpecItems(specExecutor.specification)

	beforeSpecHookStatus := specExecutor.executeBeforeSpecHook()
	if beforeSpecHookStatus.GetPassed() {
		getCurrentConsole().writeSteps(specExecutor.specification.contexts)
		dataTableRowCount := specExecutor.specification.dataTable.getRowCount()
		if dataTableRowCount == 0 {
			scenariosExecutionStatus := specExecutor.executeScenarios()
			specExecutor.specResult.addScenarioResults(scenariosExecutionStatus)
		} else {
			specExecutor.executeTableDrivenScenarios()

		}
	} else {
		addPreHook(specExecutor.specResult, beforeSpecHookStatus)
		specExecutor.currentExecutionInfo.setSpecFailure()
	}

	afterSpecHookStatus := specExecutor.executeAfterSpecHook()
	addPostHook(specExecutor.specResult, afterSpecHookStatus)

	return specExecutor.specResult
}

func (specExecutor *specExecutor) executeTableDrivenScenarios() {
	dataTableRowCount := specExecutor.specification.dataTable.getRowCount()
	for specExecutor.dataTableIndex = 0; specExecutor.dataTableIndex < dataTableRowCount; specExecutor.dataTableIndex++ {
		scenariosExecutionStatuses := specExecutor.executeScenarios()
		specExecutor.specResult.addTableDrivenScenarioResult(scenariosExecutionStatuses)
	}
}

func getTagValue(tags *tags) []string {
	tagValues := make([]string, 0)
	if tags != nil {
		tagValues = append(tagValues, tags.values...)
	}
	return tagValues
}

func (status *scenarioExecutionStatus) isFailed() bool {
	if status.hooksExecutionStatuses != nil {
		for _, hook := range status.hooksExecutionStatuses {
			if !hook.GetPassed() {
				return true
			}
		}
	}

	for _, step := range status.stepExecutionStatuses {
		if !step.passed {
			return true
		}
	}

	return false
}

func (executor *specExecutor) validateSpecification() []*stepValidationError {
	validationErrors := make([]*stepValidationError, 0)

	contextSteps := executor.specification.contexts
	contextStepsValidationErrors := executor.validateSteps(contextSteps)
	validationErrors = append(validationErrors, contextStepsValidationErrors...)

	for _, scenario := range executor.specification.scenarios {
		stepValidationErrors := executor.validateSteps(scenario.steps)
		validationErrors = append(validationErrors, stepValidationErrors...)
	}
	return validationErrors
}

func (executor *specExecutor) validateSteps(steps []*step) []*stepValidationError {
	validationErrors := make([]*stepValidationError, 0)
	for _, step := range steps {
		if step.isConcept {
			errors := executor.validateConcept(step)
			validationErrors = append(validationErrors, errors...)
		} else {
			err := executor.validateStep(step)
			if err != nil {
				validationErrors = append(validationErrors, err)
			}
		}
	}
	return validationErrors
}

func (executor *specExecutor) validateConcept(concept *step) []*stepValidationError {
	validationErrors := make([]*stepValidationError, 0)
	for _, conceptStep := range concept.conceptSteps {
		if err := executor.validateStep(conceptStep); err != nil {
			err.fileName = executor.conceptDictionary.search(concept.value).fileName
			validationErrors = append(validationErrors, err)
		}
	}
	return validationErrors
}

func (executor *specExecutor) validateStep(step *step) *stepValidationError {
	message := &Message{MessageType: Message_StepValidateRequest.Enum(),
		StepValidateRequest: &StepValidateRequest{StepText: proto.String(step.value), NumberOfArguments: proto.Int(len(step.args))}}
	response, err := getResponse(executor.connection, message)
	if err != nil {
		return &stepValidationError{step: step, message: err.Error(), fileName: executor.specification.fileName}
	}

	if response.GetMessageType() == Message_StepValidateResponse {
		validateResponse := response.GetStepValidateResponse()
		if !validateResponse.GetIsValid() {
			return &stepValidationError{step: step, message: validateResponse.GetErrorMessage(), fileName: executor.specification.fileName}
		}
	} else {
		panic("Expected a validate step response")
	}

	return nil
}

func (e *specExecutor) executeBeforeScenarioHook(scenario *scenario) *ProtoExecutionResult {
	message := &Message{MessageType: Message_ScenarioExecutionStarting.Enum(),
		ScenarioExecutionStartingRequest: &ScenarioExecutionStartingRequest{CurrentExecutionInfo: e.currentExecutionInfo}}
	e.pluginHandler.notifyPlugins(message)
	return executeAndGetStatus(e.connection, message)
}

func (executor *specExecutor) executeAfterScenarioHook() *ProtoExecutionResult {
	message := &Message{MessageType: Message_ScenarioExecutionEnding.Enum(),
		ScenarioExecutionEndingRequest: &ScenarioExecutionEndingRequest{CurrentExecutionInfo: executor.currentExecutionInfo}}
	executor.pluginHandler.notifyPlugins(message)
	return executeAndGetStatus(executor.connection, message)
}

func (specExecutor *specExecutor) executeScenarios() []*scenarioResult {
	scenarioResults := make([]*scenarioResult, 0)
	for _, scenario := range specExecutor.specification.scenarios {
		scenarioResults = append(scenarioResults, specExecutor.executeScenario(scenario))
	}
	return scenarioResults
}

func (executor *specExecutor) executeScenario(scenario *scenario) *scenarioResult {
	executor.currentExecutionInfo.CurrentScenario = &ScenarioInfo{Name: proto.String(scenario.heading.value), Tags: getTagValue(scenario.tags), IsFailed: proto.Bool(false)}
	getCurrentConsole().writeScenarioHeading(scenario.heading.value)

	scenarioResult := &scenarioResult{}
	beforeHookExecutionStatus := executor.executeBeforeScenarioHook(scenario)
	if beforeHookExecutionStatus.GetPassed() {
		contextProtoItems, passed := executor.executeContext()
		scenarioResult.addItems(contextProtoItems)
		if passed {
			scenarioProtoItems, _ := executor.executeItems(scenario.items)
			scenarioResult.addItems(scenarioProtoItems)
		}
	} else {
		addPreHook(scenarioResult, beforeHookExecutionStatus)
		executor.currentExecutionInfo.setScenarioFailure()
	}

	afterHookExecutionStatus := executor.executeAfterScenarioHook()
	addPostHook(scenarioResult, beforeHookExecutionStatus)
	return scenarioResult
}

func (executor *specExecutor) executeContext() ([]*ProtoItem, bool) {
	contextSteps := executor.specification.contexts
	items := make([]item, len(contextSteps))
	for i, context := range contextSteps {
		items[i] = context
	}
	return executor.executeItems(items)
}

func (executor *specExecutor) executeItems(items []item) ([]*ProtoItem, bool) {
	protoItems := make([]*ProtoItem, 0)
	isFailure := false
	for _, item := range items {
		protoItem := executor.executeItem(item)
		protoItems = append(protoItems, protoItem)
		if protoItem.GetItemType() == ProtoItem_Step || protoItem.GetItemType() == ProtoItem_Concept {
			if stepPassed := protoItem.GetStep().GetStepExecutionResult().ExecutionResult.Passed; !stepPassed {
				isFailure = true
				break
			}
		}
	}
	return protoItems, !isFailure
}

func (executor *specExecutor) executeItem(item item) *ProtoItem {
	if item.kind() != stepKind {
		return convertToProtoItem(item)
	}

	argLookup := new(argLookup).fromDataTableRow(&executor.specification.dataTable, executor.dataTableIndex)
	step := item.(*step)
	protoItem := &ProtoItem{}
	if step.isConcept {
		protoItem.ItemType = ProtoItem_Concept.Enum()
		protoItem.Concept = executor.executeConcept(step, argLookup)
	} else {
		protoItem.ItemType = ProtoItem_Step.Enum()
		protoItem.Step = executor.executeStep(step, argLookup)
	}
	return protoItem
}

type stepExecutionStatus struct {
	step                  *step
	resolvedArgs          []*Argument
	executionResult       []*ProtoExecutionResult
	passed                bool
	isConcept             bool
	stepExecutionStatuses []*stepExecutionStatus
}

func (s *stepExecutionStatus) addExecutionStatus(executionResult *ProtoExecutionResult) {
	if !executionResult.GetPassed() {
		s.passed = false
	}
	s.executionResult = append(s.executionResult, executionResult)
}

func (executor *specExecutor) executeSteps(steps []*step, argLookup *argLookup) []*ProtoStep {
	protoSteps := make([]*ProtoStep, 0)
	for _, step := range steps {
		protoStep := executor.executeStep(step, argLookup)
		protoSteps = append(protoSteps, protoStep)
		// TODO: handle recoverable error when verification API is done
		if !protoStep.StepExecutionResult.ExecutionResult.passed {
			break
		}
	}
	return protoSteps
}
func (executor *specExecutor) executeConcept(concept *step, dataTableLookup *argLookup) *ProtoConcept {
	conceptLookup := concept.lookup.getCopy()
	executor.populateConceptDynamicParams(conceptLookup, dataTableLookup)

	protoConcept := &ProtoConcept{ConceptStep : convertToProtoStep(concept)}
	protoConcept.Steps = executor.executeSteps(concept.conceptSteps, conceptLookup)
	for _, step := range protoConcept.Steps {
		if !step.StepExecutionResult.ExecutionResult.passed {
			conceptExecutionResult := &ProtoStepExecutionResult{ExecutionResult : step.StepExecutionResult.ExecutionResult}
			protoConcept.ConceptExecutionResult = conceptExecutionResult
			break
		}
	}
	return protoConcept

}

func printStatus(executionResult *ProtoExecutionResult) {
	getCurrentConsole().writeError(executionResult.GetErrorMessage())
	getCurrentConsole().writeError(executionResult.GetStackTrace())
}

func (executor *specExecutor) executeStep(step *step, argLookup *argLookup) *ProtoStep {

	stepRequest := executor.createStepRequest(step, argLookup)
	stepWithResolvedArgs := createStepFromStepRequest(stepRequest)
	protoStep := convertToProtoStep(stepWithResolvedArgs)
	console := getCurrentConsole()
	console.writeStep(stepWithResolvedArgs)

	protoStepExecResult := &ProtoStepExecutionResult{}
	executor.currentExecutionInfo.CurrentStep = &StepInfo{Step: stepRequest, IsFailed: proto.Bool(false)}

	beforeHookStatus := executor.executeBeforeStepHook()
	if beforeHookStatus.GetPassed() {
		executeStepMessage := &Message{MessageType: Message_ExecuteStep.Enum(), ExecuteStepRequest: stepRequest}
		stepExecutionStatus := executeAndGetStatus(executor.connection, executeStepMessage)
		if !stepExecutionStatus.GetPassed() {
			executor.currentExecutionInfo.setStepFailure()
			printStatus(stepExecutionStatus)
		}
		protoStepExecResult.ExecutionResult = stepExecutionStatus
	} else {
		protoStepExecResult.PreHookFailure = getProtoHookFailure(beforeHookStatus)
		executor.currentExecutionInfo.setStepFailure()
		printStatus(beforeHookStatus)
	}

	afterStepHookStatus := executor.executeAfterStepHook()
	if !afterStepHookStatus.GetPassed() {
		executor.currentExecutionInfo.setStepFailure()
		printStatus(afterStepHookStatus)
		protoStepExecResult.PostHookFailure = getProtoHookFailure(afterStepHookStatus)
	}

	console.writeStepFinished(stepWithResolvedArgs, *protoStepExecResult.ExecutionResult.Passed)
	protoStep.StepExecutionResult = protoStepExecResult
	return protoStep
}

func (executor *specExecutor) executeBeforeStepHook() *ProtoExecutionResult {
	message := &Message{MessageType: Message_StepExecutionStarting.Enum(),
		StepExecutionStartingRequest: &StepExecutionStartingRequest{CurrentExecutionInfo: executor.currentExecutionInfo}}
	executor.pluginHandler.notifyPlugins(message)
	return executeAndGetStatus(executor.connection, message)
}

func (executor *specExecutor) executeAfterStepHook() *ProtoExecutionResult {
	message := &Message{MessageType: Message_StepExecutionEnding.Enum(),
		StepExecutionEndingRequest: &StepExecutionEndingRequest{CurrentExecutionInfo: executor.currentExecutionInfo}}
	executor.pluginHandler.notifyPlugins(message)
	return executeAndGetStatus(executor.connection, message)
}

func populateStepExecutionResult(step *step, stepRequest *ExecuteStepRequest, result *ProtoExecutionResult) {
	status := &stepExecutionResult{isPassed: result.GetPassed(), stackTrace: result.GetStackTrace(), argument: stepRequest.GetArgs()}
	step.executionResults = append(step.executionResults, status)
}

func (executor *specExecutor) createStepRequest(step *step, argLookup *argLookup) *ExecuteStepRequest {
	stepRequest := &ExecuteStepRequest{ParsedStepText: proto.String(step.value), ActualStepText: proto.String(step.lineText)}
	stepRequest.Args = executor.createStepArgs(step.args, argLookup)
	return stepRequest
}

func (executor *specExecutor) createStepArgs(args []*stepArg, argLookup *argLookup) []*Argument {
	arguments := make([]*Argument, 0)
	for _, arg := range args {
		argument := new(Argument)
		if arg.argType == static {
			argument.Type = proto.String("string")
			argument.Value = proto.String(arg.value)
		} else if arg.argType == dynamic {
			resolvedArg := argLookup.getArg(arg.value)
			//In case a special table used in a concept, you will get a dynamic table value which has to be resolved from the concept lookup
			if resolvedArg.table.isInitialized() {
				argument.Type = proto.String("table")
				argument.Table = executor.createStepTable(&resolvedArg.table, argLookup)
			} else {
				argument.Type = proto.String("string")
				argument.Value = proto.String(resolvedArg.value)
			}
		} else {
			argument.Type = proto.String("table")
			argument.Table = executor.createStepTable(&arg.table, argLookup)
		}
		arguments = append(arguments, argument)
	}

	return arguments
}

func (executor *specExecutor) getCurrentDataTableValueFor(columnName string) string {
	return executor.specification.dataTable.get(columnName)[executor.dataTableIndex].value
}

func (executor *specExecutor) createStepTable(table *table, lookup *argLookup) *ProtoTable {
	protoTable := new(ProtoTable)
	protoTable.Headers = &ProtoTableRow{Cells: table.headers}
	tableRows := make([]*ProtoTableRow, 0)
	for i := 0; i < len(table.columns[0]); i++ {
		row := make([]string, 0)
		for _, header := range table.headers {
			tableCell := table.get(header)[i]
			value := tableCell.value
			if tableCell.cellType == dynamic {
				if lookup.containsArg(tableCell.value) {
					value = lookup.getArg(tableCell.value).value
				} else {
					//if concept has a table with dynamic cell, arglookup won't have the table value, so fetch from datatable itself
					//todo cleanup
					tableLookup := new(argLookup).fromDataTableRow(&executor.specification.dataTable, executor.dataTableIndex)
					value = tableLookup.getArg(tableCell.value).value
				}
			}
			row = append(row, value)
		}
		tableRows = append(tableRows, &ProtoTableRow{Cells: row})
	}
	protoTable.Rows = tableRows
	return protoTable
}

func executeAndGetStatus(connection net.Conn, message *Message) *ProtoExecutionResult {
	response, err := getResponse(connection, message)
	if err != nil {
		return &ProtoExecutionResult{Passed: proto.Bool(false), ErrorMessage: proto.String(err.Error())}
	}

	if response.GetMessageType() == Message_ExecutionStatusResponse {
		status := response.GetExecutionStatusResponse().GetExecutionResult()
		if status == nil {
			panic("ProtoExecutionResult should not be nil")
		}
		return status
	} else {
		panic("Expected ExecutionStatusResponse")
	}
}

func (executor *specExecutor) populateConceptDynamicParams(conceptLookup *argLookup, dataTableLookup *argLookup) {
	for key, _ := range conceptLookup.paramIndexMap {
		conceptLookupArg := conceptLookup.getArg(key)
		if conceptLookupArg.argType == dynamic {
			resolvedArg := dataTableLookup.getArg(conceptLookupArg.value)
			conceptLookup.addArgValue(key, resolvedArg)
		}
	}
}

func (executionInfo *ExecutionInfo) setSpecFailure() {
	executionInfo.CurrentSpec.IsFailed = proto.Bool(true)
}

func (executionInfo *ExecutionInfo) setScenarioFailure() {
	executionInfo.setSpecFailure()
	executionInfo.CurrentScenario.IsFailed = proto.Bool(true)
}

func (executionInfo *ExecutionInfo) setStepFailure() {
	executionInfo.setScenarioFailure()
	executionInfo.CurrentStep.IsFailed = proto.Bool(true)
}
