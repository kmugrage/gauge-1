package main

import (
	"fmt"
	"regexp"
	"strings"
)

type scenario struct {
	heading line
	steps   []*step
	tags    []string
}

type argType int

const (
	static        argType = iota
	dynamic       argType = iota
	specialString argType = iota
	specialTable  argType = iota
)

type stepArg struct {
	value   string
	argType argType
	table   table
}

type paramNameValue struct {
	name  string
	value string
}

type conceptLookup struct {
	paramIndexMap map[string]int
	paramValue    []*paramNameValue
}

type step struct {
	lineNo       int
	value        string
	lineText     string
	args         []*stepArg
	inlineTable  table
	isConcept    bool
	lookup       conceptLookup
	conceptSteps []*step
}

type specification struct {
	heading   line
	scenarios []*scenario
	comments  []*line
	dataTable table
	contexts  []*step
	fileName  string
	tags      []string
}

type line struct {
	value  string
	lineNo int
}

type parseResult struct {
	error    *parseError
	warnings []string
	ok       bool
	specFile string
}

func converterFn(predicate func(token *token, state *int) bool, apply func(token *token, spec *specification, state *int) parseResult) func(*token, *int, *specification) parseResult {

	return func(token *token, state *int, spec *specification) parseResult {
		if !predicate(token, state) {
			return parseResult{ok: true}
		}
		return apply(token, spec, state)
	}

}

func (specParser *specParser) createSpecification(tokens []*token) (*specification, *parseResult) {

	converters := initalizeConverters()
	specification := &specification{}
	finalResult := &parseResult{}
	state := initial

	for _, token := range tokens {
		for _, converter := range converters {
			result := converter(token, &state, specification)
			if !result.ok {
				if result.error != nil {
					finalResult.ok = false
					finalResult.error = result.error
					return nil, finalResult
				}
				if result.warnings != nil {
					if finalResult.warnings == nil {
						finalResult.warnings = make([]string, 0)
					}
					finalResult.warnings = append(finalResult.warnings, result.warnings...)
				}
			}
		}
	}
	finalResult.ok = true
	return specification, finalResult
}

func initalizeConverters() []func(*token, *int, *specification) parseResult {
	specConverter := converterFn(func(token *token, state *int) bool {
		return token.kind == specKind
	}, func(token *token, spec *specification, state *int) parseResult {
		if spec.heading.value != "" {
			return parseResult{ok: false, error: &parseError{token.lineNo, "Parse error: Multiple spec headings found in same file", token.lineText}}
		}
		spec.heading = line{token.value, token.lineNo}
		addStates(state, specScope)
		return parseResult{ok: true}
	})

	scenarioConverter := converterFn(func(token *token, state *int) bool {
		return token.kind == scenarioKind
	}, func(token *token, spec *specification, state *int) parseResult {
		if spec.heading.value == "" {
			return parseResult{ok: false, error: &parseError{token.lineNo, "Parse error: Scenario should be defined after the spec heading", token.lineText}}
		}
		scenarioHeading := line{token.value, token.lineNo}
		scenario := &scenario{heading: scenarioHeading}
		spec.scenarios = append(spec.scenarios, scenario)
		retainStates(state, specScope)
		addStates(state, scenarioScope)
		return parseResult{ok: true}
	})

	stepConverter := converterFn(func(token *token, state *int) bool {
		return token.kind == stepKind && isInState(*state, scenarioScope)
	}, func(token *token, spec *specification, state *int) parseResult {
		latestScenario := spec.scenarios[len(spec.scenarios)-1]
		err := spec.addStep(token, &latestScenario.steps)
		if err != nil {
			return parseResult{error: err, ok: false}
		}
		retainStates(state, specScope, scenarioScope)
		addStates(state, stepScope)
		return parseResult{ok: true}
	})

	contextConverter := converterFn(func(token *token, state *int) bool {
		return token.kind == stepKind && !isInState(*state, scenarioScope) && isInState(*state, specScope)
	}, func(token *token, spec *specification, state *int) parseResult {
		err := spec.addStep(token, &spec.contexts)
		if err != nil {
			return parseResult{error: err, ok: false}
		}
		retainStates(state, specScope)
		addStates(state, contextScope)
		return parseResult{ok: true}
	})

	commentConverter := converterFn(func(token *token, state *int) bool {
		return token.kind == commentKind
	}, func(token *token, spec *specification, state *int) parseResult {
		commentLine := &line{token.value, token.lineNo}
		spec.comments = append(spec.comments, commentLine)
		retainStates(state, specScope, scenarioScope)
		addStates(state, commentScope)
		return parseResult{ok: true}
	})

	tableHeaderConverter := converterFn(func(token *token, state *int) bool {
		return token.kind == tableHeader && isInState(*state, specScope)
	}, func(token *token, spec *specification, state *int) parseResult {
		if isInState(*state, stepScope) {
			latestScenario := spec.scenarios[len(spec.scenarios)-1]
			latestStep := latestScenario.steps[len(latestScenario.steps)-1]
			latestStep.inlineTable.addHeaders(token.args)
		} else if isInState(*state, contextScope) {
			spec.contexts[len(spec.contexts)-1].inlineTable.addHeaders(token.args)
		} else if !isInState(*state, scenarioScope) {
			if !spec.dataTable.isInitialized() {
				spec.dataTable.addHeaders(token.args)
			} else {
				value := fmt.Sprintf("multiple data table present, ignoring table at line no: %d", token.lineNo)
				return parseResult{ok: false, warnings: []string{value}}
			}
		} else {
			value := fmt.Sprintf("table not associated with a step, ignoring table at line no: %d", token.lineNo)
			return parseResult{ok: false, warnings: []string{value}}
		}
		retainStates(state, specScope, scenarioScope, stepScope, contextScope)
		addStates(state, tableScope)
		return parseResult{ok: true}
	})

	tableRowConverter := converterFn(func(token *token, state *int) bool {
		return token.kind == tableRow && isInState(*state, tableScope)
	}, func(token *token, spec *specification, state *int) parseResult {
		if isInState(*state, stepScope) {
			latestScenario := spec.scenarios[len(spec.scenarios)-1]
			latestStep := latestScenario.steps[len(latestScenario.steps)-1]
			latestStep.inlineTable.addRowValues(token.args)
		} else if isInState(*state, contextScope) {
			spec.contexts[len(spec.contexts)-1].inlineTable.addRowValues(token.args)
		} else {
			spec.dataTable.addRowValues(token.args)
		}
		retainStates(state, specScope, scenarioScope, stepScope, contextScope, tableScope)
		return parseResult{ok: true}
	})

	tagConverter := converterFn(func(token *token, state *int) bool {
		return (token.kind == tagKind)
	}, func(token *token, spec *specification, state *int) parseResult {
		if isInState(*state, scenarioScope) {
			latestScenario := spec.scenarios[len(spec.scenarios)-1]
			latestScenario.tags = token.args
		} else {
			spec.tags = token.args
		}
		return parseResult{ok: true}
	})

	converter := []func(*token, *int, *specification) parseResult{
		specConverter, scenarioConverter, stepConverter, contextConverter, commentConverter, tableHeaderConverter, tableRowConverter, tagConverter,
	}

	return converter
}

func (spec *specification) addStep(stepToken *token, addTo *[]*step) *parseError {
	step, err := spec.createStep(stepToken, false)
	if err != nil {
		return err
	}
	*addTo = append(*addTo, step)
	return nil
}

func (spec *specification) createConceptStep(token *token) (*step, *parseError) {
	return spec.createStep(token, true)
}

func (spec *specification) createStep(stepToken *token, isConcept bool) (*step, *parseError) {
	step := &step{lineNo: stepToken.lineNo, value: stepToken.value, lineText: strings.TrimSpace(stepToken.lineText)}
	r := regexp.MustCompile("{(dynamic|static|special)}")

	args := r.FindAllStringSubmatch(stepToken.value, -1)

	if args == nil {
		return step, nil
	}
	if len(args) != len(stepToken.args) {
		return nil, &parseError{stepToken.lineNo, "Step text should not have '{static}' or '{dynamic}' or '{special}'", stepToken.lineText}
	}
	step.value = r.ReplaceAllString(step.value, "{}")
	var argument *stepArg
	var err *parseError
	for i, arg := range args {
		if isConcept {
			argument, err = spec.createConceptStepArg(stepToken.args[i], arg[1], stepToken)
		} else {
			argument, err = spec.createStepArg(stepToken.args[i], arg[1], stepToken)
		}
		if err != nil {
			return nil, err
		}
		step.args = append(step.args, argument)
	}
	return step, nil
}

func (spec *specification) createStepArg(argValue string, typeOfArg string, token *token) (*stepArg, *parseError) {
	var stepArgument *stepArg
	if typeOfArg == "special" {
		return new(specialTypeResolver).resolve(argValue), nil
	} else if typeOfArg == "static" {
		return &stepArg{argType: static, value: argValue}, nil
	} else {
		if !spec.dataTable.isInitialized() {
			return nil, &parseError{token.lineNo, fmt.Sprintf("No data table found for dynamic paramter <%s> : %s", argValue, token.lineText), token.lineText}
		} else if !spec.dataTable.headerExists(argValue) {
			return nil, &parseError{token.lineNo, fmt.Sprintf("No data table column found for dynamic paramter <%s> : %s", argValue, token.lineText), token.lineText}
		}
		stepArgument = &stepArg{argType: dynamic, value: argValue}
		return stepArgument, nil
	}
}

//Does not check if data table is initialized for concepts, they will be resolved in the concept lookup
func (spec *specification) createConceptStepArg(argValue string, typeOfArg string, token *token) (*stepArg, *parseError) {
	var stepArgument *stepArg
	if typeOfArg == "special" {
		return new(specialTypeResolver).resolve(argValue), nil
	} else if typeOfArg == "static" {
		return &stepArg{argType: static, value: argValue}, nil
	} else {
		stepArgument = &stepArg{argType: dynamic, value: argValue}
		return stepArgument, nil
	}

}

type specialTypeResolver struct {
}

func (resolver *specialTypeResolver) resolve(value string) *stepArg {
	return &stepArg{argType: specialString, value: ""}
}

func (lookup *conceptLookup) addParam(param string) {
	if lookup.paramIndexMap == nil {
		lookup.paramIndexMap = make(map[string]int)
		lookup.paramValue = make([]*paramNameValue, 0)
	}
	lookup.paramIndexMap[param] = len(lookup.paramValue)
	lookup.paramValue = append(lookup.paramValue, &paramNameValue{name: param})
}

func (lookup *conceptLookup) containsParam(param string) bool {
	index := lookup.paramIndexMap[param]
	return lookup.paramValue[index].name == param
}
