package main

import (
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
	oldStep, _ := spec.createStep(stepTokens[0])
	newStep, _ := spec.createStep(stepTokens[1])

	r := &refactorer{oldStep: oldStep, newStep: newStep}
	return r, nil
}

func (r *refactorer) performRefactoring() error {
	fmt.Println(r.oldStep.lineText)
	for _, f := range r.oldStep.fragments {
		fmt.Printf("%d %s\n", f.GetFragmentType(), f.GetText())
	}

	fmt.Println(r.newStep.lineText)
	for _, f := range r.newStep.fragments {
		fmt.Printf("%d %s\n", f.GetFragmentType(), f.GetText())
	}
	return nil
}
