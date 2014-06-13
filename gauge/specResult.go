package main

import (
	"code.google.com/p/goprotobuf/proto"
)

type suiteResult struct {
	specResults          []*specResult
	preSuite             *ProtoHookFailure
	postSuite            *ProtoHookFailure
	currentSpecIndex     int
	currentScenarioIndex int
}

type specResult struct {
	preSpec         *ProtoHookFailure
	protoSpec        *ProtoSpec
	postSpec        *ProtoHookFailure
}

type scenarioResult struct {
	preScenario         *ProtoHookFailure
	protoScenario       *ProtoScenario
	postScenario        *ProtoHookFailure
}

type result interface {
	getPreHook() **ProtoHookFailure
	getPostHook() **ProtoHookFailure
}

func (suiteResult *suiteResult) getPreHook() **ProtoHookFailure {
	return &suiteResult.preSuite
}

func (suiteResult *suiteResult) getPostHook() **ProtoHookFailure {
	return &suiteResult.postSuite
}

func (specResult *specResult) getPreHook() **ProtoHookFailure {
	return &specResult.preSpec
}

func (specResult *specResult) getPostHook() **ProtoHookFailure {
	return &specResult.postSpec
}

func (scenarioResult *scenarioResult) getPreHook() **ProtoHookFailure {
	return &scenarioResult.preScenario
}

func (scenarioResult *scenarioResult) getPostHook() **ProtoHookFailure {
	return &scenarioResult.postScenario
}

func (specResult *specResult) addSpecItems(spec *specification) {
	for _, item := range spec.items {
		if item.kind() != scenarioKind {
			specResult.protoSpec.Items = append(specResult.protoSpec.Items, convertToProtoItem(item))
		}
	}
}

func newSuiteResult() *suiteResult {
	result := new(suiteResult)
	result.specResults = make([]*specResult, 0)
	result.currentSpecIndex = -1
	return result
}

func addPreHook(result result, executionResult *ProtoExecutionResult) {
	if !executionResult.GetPassed() {
		*(result.getPreHook()) = getProtoHookFailure(executionResult)
	}
}

func addPostHook(result result, executionResult *ProtoExecutionResult) {
	if !executionResult.GetPassed() {
		*(result.getPostHook()) = getProtoHookFailure(executionResult)
	}
}

func (suiteResult *suiteResult) addSpecResult(specResult *specResult) {
	suiteResult.specResults = append(suiteResult.specResults, specResult)
	suiteResult.currentSpecIndex++
}

func getProtoHookFailure(executionResult *ProtoExecutionResult) *ProtoHookFailure {
	return &ProtoHookFailure{StackTrace: executionResult.StackTrace, ErrorMessage: executionResult.ErrorMessage, ScreenShot: executionResult.ScreenShot}
}

func (suiteResult *suiteResult) startTableDrivenScenarios() {
	suiteResult.getCurrentSpec().IsTableDriven = proto.Bool(true)
}

func (suiteResult *suiteResult) getCurrentSpec() *ProtoSpec {
	return suiteResult.protoSpecResult[suiteResult.currentSpecIndex]
}

func (suiteResult *suiteResult) newSpecStart() {
	suiteResult.currentSpecIndex++
	suiteResult.currentScenarioIndex = -1
	suiteResult.protoSpecResult = append(suiteResult.protoSpecResult, new(ProtoSpec))
}

func (suiteResult *suiteResult) newScenarioStart() {
	suiteResult.currentScenarioIndex++
	suiteResult.getCurrentSpec().Items = append(suiteResult.getCurrentSpec().Items, new(ProtoItem))
}
func (specResult *specResult) addScenarioResults(scenarioResult []*scenarioResult) []*scenarioResult {
	if (specResult.protoSpec == nil ) {
		specResult.protoSpec = &ProtoSpec{Items:make([]*ProtoItem, 0)}
	}
	specResult.protoSpec.Items = append(specResult.protoSpec.Items, &ProtoScenario{})
}

func (scenarioResult *scenarioResult) addItems(protoItems []*ProtoItem) {
	if (scenarioResult.protoScenario == nil) {
		scenarioResult.protoScenario = &ProtoScenario{ScenarioItems:make([]*ProtoItem, 0)}
	}
	scenarioResult.protoScenario.ScenarioItems = append(scenarioResult.protoScenario.ScenarioItems, protoItems)
}
