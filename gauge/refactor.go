package main

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
)

type refactorer struct {
	oldStep *step
	newStep *step
}

func newRefactorer(oldStepText, newStepText string) (*refactorer, error) {
	parser := new(specParser)
	stepTokens, err := parser.generateTokens("* " + oldStepText + "\n" + "*" + newStepText)
	if err != nil {
		return nil, err
	}

	spec := &specification{}
	oldStep, err := spec.createStepUsingLookup(stepTokens[0], nil)
	if err != nil {
		return nil, err
	}
	newStep, err := spec.createStepUsingLookup(stepTokens[1], nil)
	if err != nil {
		return nil, err
	}

	r := &refactorer{oldStep: oldStep, newStep: newStep}
	return r, nil
}

func (r *refactorer) performRefactoring() error {
	request := &RefactorRequest{OldStepText: proto.String(r.oldStep.value), NewStepText: proto.String(r.newStep.value)}
	fmt.Println(request.GetNewStepText())
	return nil
}
