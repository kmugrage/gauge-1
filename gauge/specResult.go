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
	protoSpecResult *ProtoSpec
	postSpec        *ProtoHookFailure
}

type scenarioResult struct {
	preScenario         *ProtoHookFailure
	protoScenarioResult *ProtoScenario
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
			specResult.protoSpecResult.Items = append(specResult.protoSpecResult.Items, convertToProtoItem(item))
		}
	}
}

func newSuiteResult() *suiteResult {
	result := new(suiteResult)
	result.specResults = make([]*specResult, 0)
	result.currentSpecIndex = -1
	return result
}

func addPreHook(result result, execStatus *ExecutionStatus) {
	if !execStatus.GetPassed() {
		*(result.getPreHook()) = getProtoHookFailure(execStatus)
	}
}

func addPostHook(result result, execStatus *ExecutionStatus) {
	if !execStatus.GetPassed() {
		*(result.getPostHook()) = getProtoHookFailure(execStatus)
	}
}

func (suiteResult *suiteResult) addSpecResult(specResult *specResult) {
	suiteResult.specResults = append(suiteResult.specResults, specResult)
	suiteResult.currentSpecIndex++
}

func getProtoHookFailure(execStatus *ExecutionStatus) *ProtoHookFailure {
	return &ProtoHookFailure{StackTrace: execStatus.StackTrace, ErrorMessage: execStatus.ErrorMessage, ScreenShot: execStatus.ScreenShot}
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
