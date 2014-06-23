package main

type refactorer struct {
	oldStep *step
	newStep *step
}

func newRefactorer(oldStepText, newStepText string) (*refactorer, error) {
	parser := new(specParser)
	oldStepTokens, err := parser.generateTokens("* " + oldStepText)
	if err != nil {
		return nil, err
	}
	newStepTokens, err := parser.generateTokens("* " + newStepText)
	if err != nil {
		return nil, err
	}
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: scenarioKind, value: "Scenario heading", lineNo: 2},
	}

	dummySpec, _ := new(specParser).createSpecification(tokens, new(conceptDictionary))
	oldStep, _ := dummySpec.createStep(oldStepTokens[0])
	newStep, _ := dummySpec.createStep(newStepTokens[0])

	r := &refactorer{oldStep: oldStep, newStep: newStep}
	return r, nil
}

func (r *refactorer) performRefactoring() error {
	return nil
}
