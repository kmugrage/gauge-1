package main

import "fmt"

type conceptDictionary struct {
	conceptsMap map[string]*step
}

type conceptParser struct {
	currentState   int
	currentConcept *step
}

func (parser *conceptParser) parse(text string) ([]*step, *parseError) {
	defer parser.resetState()

	specParser := new(specParser)
	tokens, err := specParser.generateTokens(text)
	if err != nil {
		return nil, err
	}
	return parser.createConcepts(tokens)
}

func (parser *conceptParser) resetState() {
	parser.currentState = 0
	parser.currentConcept = nil
}

func (parser *conceptParser) createConcepts(tokens []*token) ([]*step, *parseError) {
	concepts := make([]*step, 0)
	var error *parseError
	for _, token := range tokens {
		if parser.isConceptHeading(token) {
			if isInState(parser.currentState, conceptScope, stepScope) {
				concepts = append(concepts, parser.currentConcept)
			}
			addStates(&parser.currentState, conceptScope)
			parser.currentConcept, error = parser.processConceptHeading(token)
			if error != nil {
				return nil, error
			}
		} else if parser.isStep(token) {
			if !isInState(parser.currentState, conceptScope) {
				return nil, &parseError{lineNo: token.lineNo, message: "Step is not defined inside a concept heading", lineText: token.lineText}
			}
			if err := parser.processConceptStep(token); err != nil {
				return nil, err
			}
			addStates(&parser.currentState, stepScope)
		}
	}
	if !isInState(parser.currentState, stepScope) {
		return nil, &parseError{lineNo: parser.currentConcept.lineNo, message: "Concept should have atleast one step", lineText: parser.currentConcept.lineText}
	}
	return append(concepts, parser.currentConcept), nil
}

func (parser *conceptParser) isConceptHeading(token *token) bool {
	return token.kind == specKind || token.kind == scenarioKind
}

func (parser *conceptParser) isStep(token *token) bool {
	return token.kind == stepKind
}

func (parser *conceptParser) processConceptHeading(token *token) (*step, *parseError) {
	processStep(new(specParser), token)
	concept, err := new(specification).createConceptStep(token)
	if err != nil {
		return nil, err
	}
	concept.isConcept = true
	if !parser.hasOnlyDynamicParams(concept) {
		return nil, &parseError{lineNo: concept.lineNo, message: "Concept heading can have only Dynamic Parameters"}
	}
	parser.createConceptLookup(concept)
	return concept, nil

}

func (parser *conceptParser) processConceptStep(token *token) *parseError {
	processStep(new(specParser), token)
	conceptStep, err := new(specification).createConceptStep(token)
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

func (parser *conceptParser) validateConceptStep(conceptStep *step) *parseError {
	for _, arg := range conceptStep.args {
		if arg.argType == dynamic && !parser.currentConcept.lookup.containsParam(arg.value) {
			return &parseError{lineNo: conceptStep.lineNo, message: fmt.Sprintf("Dynamic parameter <%s> is not defined in concept heading", arg.value)}
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
