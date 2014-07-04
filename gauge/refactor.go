package main

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"net"
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

func (r *refactorer) performRefactoring(conn net.Conn) error {
	request := &RefactorRequest{OldStepText: proto.String(r.oldStep.value), NewStepText: proto.String(r.newStep.value)}
	for _, fragment := range r.newStep.fragments {
		if fragment.GetFragmentType() == Fragment_Parameter {
			request.Params = append(request.Params, fragment.Parameter)
		}
	}

	message := &Message{MessageType: Message_RefactorRequest.Enum(), RefactorRequest: request}
	_, err := getResponse(conn, message)
	if err != nil {
		return err
	}

	fmt.Println(request.GetNewStepText())
	fmt.Printf("%v\n", request)
	return nil
}
