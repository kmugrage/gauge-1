package main

import (
	"code.google.com/p/goprotobuf/proto"
)

type suiteResult struct {
	specResults           []*specResult
	preSuite             *ProtoHookFailure
	postSuite            *ProtoHookFailure
	currentSpecIndex      int
	currentScenarioIndex  int
}

type specResult struct {
	preSuite             *ProtoHookFailure
	protoSpecResult      *ProtoSpec
	postSuite            *ProtoHookFailure
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
	return &specResult.preSuite
}

func (specResult *specResult) getPostHook() **ProtoHookFailure {
	return &specResult.preSuite
}

func newSuiteResult() *suiteResult {
	result := new(suiteResult)
	result.specResults = make([]*specResult, 0)
	result.currentSpecIndex = -1
	return result
}

func addPreHook(result result, execStatus *ExecutionStatus) {
	if !execStatus.GetPassed() {
		*(result.getPreHook()) = &ProtoHookFailure{StackTrace:execStatus.StackTrace, ErrorMessage:execStatus.ErrorMessage, ScreenShot:execStatus.ScreenShot}
	}
}

func addPostHook(result result, execStatus *ExecutionStatus) {
	if !execStatus.GetPassed() {
		*(result.getPostHook()) = &ProtoHookFailure{StackTrace:execStatus.StackTrace, ErrorMessage:execStatus.ErrorMessage, ScreenShot:execStatus.ScreenShot}
	}
}

func (suiteResult *suiteResult) addSpecResult(specResult *specResult) {
	suiteResult.specResults = append(suiteResult.specResults, specResult)
	suiteResult.currentSpecIndex ++
}

func (suiteResult *suiteResult) addPreHook(execStatus *ExecutionStatus) {
	if !execStatus.GetPassed() {
		suiteResult.preSuite = &ProtoHookFailure{StackTrace:execStatus.StackTrace, ErrorMessage:execStatus.ErrorMessage, ScreenShot:execStatus.ScreenShot}
	}
}

func (suiteResult *suiteResult) addPostHook(execStatus *ExecutionStatus) {
	if !execStatus.GetPassed() {
		suiteResult.postSuite = &ProtoHookFailure{StackTrace:execStatus.StackTrace, ErrorMessage:execStatus.ErrorMessage, ScreenShot:execStatus.ScreenShot}
	}
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
