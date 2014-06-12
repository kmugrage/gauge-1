package main

import (
	"code.google.com/p/goprotobuf/proto"
)

type suiteResult struct {
	protoSpecResult      []*ProtoSpec
	preSuite             ProtoHookFailure
	postSuite            ProtoHookFailure
	currentSpecIndex     int
	currentScenarioIndex int
}

func newSuiteResult() *suiteResult {
	result := new(suiteResult)
	result.protoSpecResult = make([]*ProtoSpec, 0)
	result.currentSpecIndex = -1
	return result
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
