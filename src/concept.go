package main

import "fmt"

type conceptDictionary struct {
	conceptsMap map[string]*step
}

type conceptParser struct {
	currentState   int
	currentConcept *step
}

func (parser *conceptParser) parse(text string) ([]*step, error) {
	specParser := new(specParser)
	tokens, err := specParser.generateTokens(text)
	if err != nil {
		return nil, err
	}
	return parser.createConcepts(tokens)
}

func (parser *conceptParser) createConcepts(tokens []*token) ([]*step, error) {
	concepts := make([]*step, 0)
	var err error
	for _, token := range tokens {
		if parser.isConceptHeading(token) {
			if isInState(parser.currentState, conceptScope, stepScope) {
				concepts = append(concepts, parser.currentConcept)
			}
			addStates(&parser.currentState, conceptScope)
			parser.currentConcept, err = parser.processConceptHeading(token)
			if err != nil {
				return nil, err
			}
		} else if parser.isStep(token) {
			if !isInState(parser.currentState, conceptScope) {
				return nil, &syntaxError{lineNo: token.lineNo, message: "Step is not defined inside a concept heading"}
			}
			if err := parser.processConceptStep(token); err != nil {
				return nil, err
			}
			addStates(&parser.currentState, stepScope)
		}
	}
	if !isInState(parser.currentState, stepScope) {
		return nil, &syntaxError{lineNo: parser.currentConcept.lineNo, message: "Concept should have atleast one step"}
	}
	return append(concepts, parser.currentConcept), nil
}

func (parser *conceptParser) isConceptHeading(token *token) bool {
	return token.kind == specKind || token.kind == scenarioKind
}

func (parser *conceptParser) isStep(token *token) bool {
	return token.kind == stepKind
}

func (parser *conceptParser) processConceptHeading(token *token) (*step, error) {
	processStep(new(specParser), token)
	concept, err := new(specification).createStep(token)
	if err != nil {
		return nil, err
	}
	concept.isConcept = true
	if !parser.hasOnlyDynamicParams(concept) {
		return nil, &syntaxError{lineNo: concept.lineNo, message: "Concept heading can have only Dynamic Parameters"}
	}
	parser.createConceptLookup(concept)
	return concept, nil

}

func (parser *conceptParser) processConceptStep(token *token) error {
	processStep(new(specParser), token)
	conceptStep, err := new(specification).createStep(token)
	if err != nil {
		return err
	}
	if err := parser.validateConceptStep(conceptStep); err != nil {
		return err
	}
	parser.currentConcept.conceptSteps = append(parser.currentConcept.conceptSteps, conceptStep)
	return nil
}

func (parser *conceptParser) hasOnlyDynamicParams(concept *step) bool {
	for _, arg := range concept.args {
		if arg.argType != dynamic {
			return false
		}
	}
	return true
}

func (parser *conceptParser) validateConceptStep(conceptStep *step) error {
	for _, arg := range conceptStep.args {
		if arg.argType == dynamic && !parser.currentConcept.lookup.containsParam(arg.value) {
			return &syntaxError{lineNo: conceptStep.lineNo, message: fmt.Sprintf("dynamic parameter %s is not defined in concept heading", arg.value)}
		}
	}
	return nil
}

func (parser *conceptParser) createConceptLookup(concept *step) {
	for _, arg := range concept.args {
		concept.lookup.addParam(arg.value)
	}
}

func (conceptDictionary *conceptDictionary) add(concepts []*step) {
	for _, step := range concepts {
		conceptDictionary.conceptsMap[step.value] = step
	}
}
