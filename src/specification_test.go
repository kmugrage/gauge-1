package main

import . "launchpad.net/gocheck"

func (s *MySuite) TestThrowsErrorForMultipleSpecHeading(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 2},
		&token{kind: stepKind, value: "Example step", lineNo: 3},
		&token{kind: specKind, value: "Another Heading", lineNo: 4},
	}

	_, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, false)

	c.Assert(result.error.message, Equals, "Parse error: Multiple spec headings found in same file")
	c.Assert(result.error.lineNo, Equals, 4)
}

func (s *MySuite) TestThrowsErrorForScenarioWithoutSpecHeading(c *C) {
	tokens := []*token{
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 1},
		&token{kind: stepKind, value: "Example step", lineNo: 2},
	}

	_, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, false)

	c.Assert(result.error.message, Equals, "Parse error: Scenario should be defined after the spec heading")
	c.Assert(result.error.lineNo, Equals, 1)
}

func (s *MySuite) TestSpecWithHeadingAndSimpleSteps(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 2},
		&token{kind: stepKind, value: "Example step", lineNo: 3},
	}

	spec, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, true)
	c.Assert(spec.heading.lineNo, Equals, 1)
	c.Assert(spec.heading.value, Equals, "Spec Heading")

	c.Assert(len(spec.scenarios), Equals, 1)
	c.Assert(spec.scenarios[0].heading.lineNo, Equals, 2)
	c.Assert(spec.scenarios[0].heading.value, Equals, "Scenario Heading")
	c.Assert(len(spec.scenarios[0].steps), Equals, 1)
	c.Assert(spec.scenarios[0].steps[0].value, Equals, "Example step")
}

func (s *MySuite) TestStepsAndComments(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: commentKind, value: "A comment with some text and **bold** characters", lineNo: 2},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 3},
		&token{kind: commentKind, value: "Another comment", lineNo: 4},
		&token{kind: stepKind, value: "Example step", lineNo: 5},
		&token{kind: commentKind, value: "Third comment", lineNo: 6},
	}

	spec, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, true)
	c.Assert(spec.heading.value, Equals, "Spec Heading")

	c.Assert(len(spec.comments), Equals, 3)
	c.Assert(spec.comments[0].lineNo, Equals, 2)
	c.Assert(spec.comments[0].value, Equals, "A comment with some text and **bold** characters")

	c.Assert(len(spec.scenarios), Equals, 1)
	scenario := spec.scenarios[0]

	c.Assert(spec.comments[1].lineNo, Equals, 4)
	c.Assert(spec.comments[1].value, Equals, "Another comment")

	c.Assert(spec.comments[2].lineNo, Equals, 6)
	c.Assert(spec.comments[2].value, Equals, "Third comment")

	c.Assert(scenario.heading.value, Equals, "Scenario Heading")
	c.Assert(len(scenario.steps), Equals, 1)
}

func (s *MySuite) TestStepsWithParam(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: tableHeader, args: []string{"id"}, lineNo: 2},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 3},
		&token{kind: stepKind, value: "enter {static} with {dynamic}", lineNo: 4, args: []string{"user", "id"}},
		&token{kind: stepKind, value: "sample \\{static\\}", lineNo: 5},
	}

	spec, result := new(specParser).createSpecification(tokens, new(conceptDictionary))
	c.Assert(result.ok, Equals, true)
	step := spec.scenarios[0].steps[0]
	c.Assert(step.value, Equals, "enter {} with {}")
	c.Assert(step.lineNo, Equals, 4)
	c.Assert(len(step.args), Equals, 2)
	c.Assert(step.args[0].value, Equals, "user")
	c.Assert(step.args[0].argType, Equals, static)
	c.Assert(step.args[1].value, Equals, "id")
	c.Assert(step.args[1].argType, Equals, dynamic)

	escapedStep := spec.scenarios[0].steps[1]
	c.Assert(escapedStep.value, Equals, "sample \\{static\\}")
	c.Assert(len(escapedStep.args), Equals, 0)
}

func (s *MySuite) TestStepsWithKeywords(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 2},
		&token{kind: stepKind, value: "sample {static} and {dynamic}", lineNo: 3, args: []string{"name"}},
	}

	_, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result, NotNil)
	c.Assert(result.ok, Equals, false)
	c.Assert(result.error.message, Equals, "Step text should not have '{static}' or '{dynamic}' or '{special}'")
	c.Assert(result.error.lineNo, Equals, 3)
}

func (s *MySuite) TestContextWithKeywords(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: stepKind, value: "sample {static} and {dynamic}", lineNo: 3, args: []string{"name"}},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 2},
	}

	_, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result, NotNil)
	c.Assert(result.ok, Equals, false)
	c.Assert(result.error.message, Equals, "Step text should not have '{static}' or '{dynamic}' or '{special}'")
	c.Assert(result.error.lineNo, Equals, 3)
}

func (s *MySuite) TestSpecWithDataTable(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading"},
		&token{kind: commentKind, value: "Comment before data table"},
		&token{kind: tableHeader, args: []string{"id", "name"}},
		&token{kind: tableRow, args: []string{"1", "foo"}},
		&token{kind: tableRow, args: []string{"2", "bar"}},
		&token{kind: commentKind, value: "Comment before data table"},
	}

	spec, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, true)
	c.Assert(spec.dataTable, NotNil)
	c.Assert(len(spec.dataTable.get("id")), Equals, 2)
	c.Assert(len(spec.dataTable.get("name")), Equals, 2)
	c.Assert(spec.dataTable.get("id")[0], Equals, "1")
	c.Assert(spec.dataTable.get("id")[1], Equals, "2")
	c.Assert(spec.dataTable.get("name")[0], Equals, "foo")
	c.Assert(spec.dataTable.get("name")[1], Equals, "bar")
}

func (s *MySuite) TestStepWithInlineTable(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 2},
		&token{kind: stepKind, value: "Step with inline table", lineNo: 3},
		&token{kind: tableHeader, args: []string{"id", "name"}},
		&token{kind: tableRow, args: []string{"1", "foo"}},
		&token{kind: tableRow, args: []string{"2", "bar"}},
	}

	spec, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, true)
	inlineTable := spec.scenarios[0].steps[0].inlineTable
	c.Assert(inlineTable, NotNil)
	c.Assert(len(inlineTable.get("id")), Equals, 2)
	c.Assert(len(inlineTable.get("name")), Equals, 2)
	c.Assert(inlineTable.get("id")[0], Equals, "1")
	c.Assert(inlineTable.get("id")[1], Equals, "2")
	c.Assert(inlineTable.get("name")[0], Equals, "foo")
	c.Assert(inlineTable.get("name")[1], Equals, "bar")
}

func (s *MySuite) TestContextWithInlineTable(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading"},
		&token{kind: stepKind, value: "Context with inline table"},
		&token{kind: tableHeader, args: []string{"id", "name"}},
		&token{kind: tableRow, args: []string{"1", "foo"}},
		&token{kind: tableRow, args: []string{"2", "bar"}},
		&token{kind: scenarioKind, value: "Scenario Heading"},
	}

	spec, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, true)
	inlineTable := spec.contexts[0].inlineTable

	c.Assert(inlineTable, NotNil)
	c.Assert(len(inlineTable.get("id")), Equals, 2)
	c.Assert(len(inlineTable.get("name")), Equals, 2)
	c.Assert(inlineTable.get("id")[0], Equals, "1")
	c.Assert(inlineTable.get("id")[1], Equals, "2")
	c.Assert(inlineTable.get("name")[0], Equals, "foo")
	c.Assert(inlineTable.get("name")[1], Equals, "bar")
}

func (s *MySuite) TestWarningWhenParsingMultipleDataTable(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading"},
		&token{kind: commentKind, value: "Comment before data table"},
		&token{kind: tableHeader, args: []string{"id", "name"}},
		&token{kind: tableRow, args: []string{"1", "foo"}},
		&token{kind: tableRow, args: []string{"2", "bar"}},
		&token{kind: commentKind, value: "Comment before data table"},
		&token{kind: tableHeader, args: []string{"phone"}, lineNo: 7},
		&token{kind: tableRow, args: []string{"1"}},
		&token{kind: tableRow, args: []string{"2"}},
	}

	_, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, true)
	c.Assert(len(result.warnings), Equals, 1)
	c.Assert(result.warnings[0], Equals, "multiple data table present, ignoring table at line no: 7")

}

func (s *MySuite) TestWarningWhenParsingTableOccursWithoutStep(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 2},
		&token{kind: tableHeader, args: []string{"id", "name"}, lineNo: 3},
		&token{kind: tableRow, args: []string{"1", "foo"}, lineNo: 4},
		&token{kind: tableRow, args: []string{"2", "bar"}, lineNo: 5},
		&token{kind: stepKind, value: "Step", lineNo: 6},
		&token{kind: commentKind, value: "comment in between", lineNo: 7},
		&token{kind: tableHeader, args: []string{"phone"}, lineNo: 8},
		&token{kind: tableRow, args: []string{"1"}},
		&token{kind: tableRow, args: []string{"2"}},
	}

	_, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, true)
	c.Assert(len(result.warnings), Equals, 2)
	c.Assert(result.warnings[0], Equals, "table not associated with a step, ignoring table at line no: 3")
	c.Assert(result.warnings[1], Equals, "table not associated with a step, ignoring table at line no: 8")

}

func (s *MySuite) TestAddSpecTags(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: tagKind, args: []string{"tag1", "tag2"}, lineNo: 2},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 3},
	}

	spec, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, true)

	c.Assert(len(spec.tags), Equals, 2)
	c.Assert(spec.tags[0], Equals, "tag1")
	c.Assert(spec.tags[1], Equals, "tag2")
}

func (s *MySuite) TestAddSpecTagsAndScenarioTags(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: tagKind, args: []string{"tag1", "tag2"}, lineNo: 2},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 3},
		&token{kind: tagKind, args: []string{"tag3", "tag4"}, lineNo: 2},
	}

	spec, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, true)

	c.Assert(len(spec.tags), Equals, 2)
	c.Assert(spec.tags[0], Equals, "tag1")
	c.Assert(spec.tags[1], Equals, "tag2")

	c.Assert(len(spec.scenarios[0].tags), Equals, 2)
	c.Assert(spec.scenarios[0].tags[0], Equals, "tag3")
	c.Assert(spec.scenarios[0].tags[1], Equals, "tag4")
}

func (s *MySuite) TestErrorOnAddingDynamicParamterWithoutADataTable(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 2},
		&token{kind: stepKind, value: "Step with a {dynamic}", args: []string{"foo"}, lineNo: 3, lineText: "*Step with a <foo>"},
	}

	_, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, false)
	c.Assert(result.error.message, Equals, "Dynamic parameter <foo> could not be resolved")
	c.Assert(result.error.lineNo, Equals, 3)

}

func (s *MySuite) TestErrorOnAddingDynamicParamterWithoutDataTableHeaderValue(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: tableHeader, args: []string{"id, name"}, lineNo: 2},
		&token{kind: tableRow, args: []string{"123, hello"}, lineNo: 3},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 4},
		&token{kind: stepKind, value: "Step with a {dynamic}", args: []string{"foo"}, lineNo: 5, lineText: "*Step with a <foo>"},
	}

	_, result := new(specParser).createSpecification(tokens, new(conceptDictionary))

	c.Assert(result.ok, Equals, false)
	c.Assert(result.error.message, Equals, "Dynamic parameter <foo> could not be resolved")
	c.Assert(result.error.lineNo, Equals, 5)

}

func (s *MySuite) TestLookupAddParam(c *C) {
	lookup := new(conceptLookup)
	lookup.addParam("param1")
	lookup.addParam("param2")

	c.Assert(lookup.paramIndexMap["param1"], Equals, 0)
	c.Assert(lookup.paramIndexMap["param2"], Equals, 1)
	c.Assert(len(lookup.paramValue), Equals, 2)
	c.Assert(lookup.paramValue[0].name, Equals, "param1")
	c.Assert(lookup.paramValue[1].name, Equals, "param2")

}

func (s *MySuite) TestLookupContainsParam(c *C) {
	lookup := new(conceptLookup)
	lookup.addParam("param1")
	lookup.addParam("param2")

	c.Assert(lookup.containsParam("param1"), Equals, true)
	c.Assert(lookup.containsParam("param2"), Equals, true)
	c.Assert(lookup.containsParam("param3"), Equals, false)
}

func (s *MySuite) TestAddParamValue(c *C) {
	lookup := new(conceptLookup)
	lookup.addParam("param1")
	lookup.addParamValue("param1", "value1", static)
	lookup.addParam("param2")
	lookup.addParamValue("param2", "value2", dynamic)

	c.Assert(lookup.getParamValue("param1"), Equals, "value1")
	c.Assert(lookup.getParamValue("param2"), Equals, "value2")
}

func (s *MySuite) TestPanicForInvalidParam(c *C) {
	lookup := new(conceptLookup)

	c.Assert(func() { lookup.addParamValue("param1", "value1", static) }, Panics, "Accessing an invalid parameter (param1)")
	c.Assert(func() { lookup.getParamValue("param1") }, Panics, "Accessing an invalid parameter (param1)")
}

func (s *MySuite) TestGetLookupCopy(c *C) {
	originalLookup := new(conceptLookup)
	originalLookup.addParam("param1")

	copiedLookup := originalLookup.getCopy()
	copiedLookup.addParamValue("param1", "value1", static)

	c.Assert(copiedLookup.getParamValue("param1"), Equals, "value1")
	c.Assert(originalLookup.getParamValue("param1"), Equals, "")
}

func (s *MySuite) TestCreateStepFromSimpleConcept(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 2},
		&token{kind: stepKind, value: "concept step", lineNo: 3},
	}

	conceptDictionary := new(conceptDictionary)
	firstStep := &step{value: "step 1"}
	secondStep := &step{value: "step 2"}
	conceptStep := &step{value: "concept step", isConcept: true, conceptSteps: []*step{firstStep, secondStep}}
	conceptDictionary.add([]*step{conceptStep})
	spec, result := new(specParser).createSpecification(tokens, conceptDictionary)
	c.Assert(result.ok, Equals, true)

	c.Assert(len(spec.scenarios[0].steps), Equals, 1)
	specConceptStep := spec.scenarios[0].steps[0]
	c.Assert(specConceptStep.isConcept, Equals, true)
	c.Assert(specConceptStep.conceptSteps[0], Equals, firstStep)
	c.Assert(specConceptStep.conceptSteps[1], Equals, secondStep)
}

func (s *MySuite) TestCreateStepFromConceptWithParameters(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 2},
		&token{kind: stepKind, value: "create user {static}", args: []string{"foo"}, lineNo: 3},
		&token{kind: stepKind, value: "create user {static}", args: []string{"bar"}, lineNo: 4},
	}

	concepts, _ := new(conceptParser).parse("#create user <username> \n * enter user <username> \n *select \"finish\"")
	conceptsDictionary := new(conceptDictionary)
	conceptsDictionary.add(concepts)

	spec, result := new(specParser).createSpecification(tokens, conceptsDictionary)
	c.Assert(result.ok, Equals, true)

	c.Assert(len(spec.scenarios[0].steps), Equals, 2)

	firstConceptStep := spec.scenarios[0].steps[0]
	c.Assert(firstConceptStep.isConcept, Equals, true)
	c.Assert(firstConceptStep.conceptSteps[0].value, Equals, "enter user {}")
	c.Assert(firstConceptStep.conceptSteps[0].args[0].value, Equals, "username")
	c.Assert(firstConceptStep.conceptSteps[1].value, Equals, "select {}")
	c.Assert(firstConceptStep.conceptSteps[1].args[0].value, Equals, "finish")
	c.Assert(len(firstConceptStep.lookup.paramValue), Equals, 1)
	c.Assert(firstConceptStep.lookup.paramValue[0].name, Equals, "username")
	c.Assert(firstConceptStep.lookup.paramValue[0].value, Equals, "foo")

	secondConceptStep := spec.scenarios[0].steps[1]
	c.Assert(secondConceptStep.isConcept, Equals, true)
	c.Assert(secondConceptStep.conceptSteps[0].value, Equals, "enter user {}")
	c.Assert(secondConceptStep.conceptSteps[0].args[0].value, Equals, "username")
	c.Assert(secondConceptStep.conceptSteps[1].value, Equals, "select {}")
	c.Assert(secondConceptStep.conceptSteps[1].args[0].value, Equals, "finish")
	c.Assert(len(secondConceptStep.lookup.paramValue), Equals, 1)
	c.Assert(secondConceptStep.lookup.paramValue[0].name, Equals, "username")
	c.Assert(secondConceptStep.lookup.paramValue[0].value, Equals, "bar")

}

func (s *MySuite) TestCreateStepFromConceptWithDynamicParameters(c *C) {
	tokens := []*token{
		&token{kind: specKind, value: "Spec Heading", lineNo: 1},
		&token{kind: tableHeader, args: []string{"id, description"}, lineNo: 2},
		&token{kind: tableRow, args: []string{"123, Admin fellow"}, lineNo: 3},
		&token{kind: scenarioKind, value: "Scenario Heading", lineNo: 4},
		&token{kind: stepKind, value: "create user {dynamic} and {dynamic}", args: []string{"id", "description"}, lineNo: 5},
		&token{kind: stepKind, value: "create user {static} and {static}", args: []string{"456", "Regular fellow"}, lineNo: 6},
	}

	concepts, _ := new(conceptParser).parse("#create user <user-id> and <user-description> \n * enter user <user-id> and <user-description> \n *select \"finish\"")
	conceptsDictionary := new(conceptDictionary)
	conceptsDictionary.add(concepts)

	spec, result := new(specParser).createSpecification(tokens, conceptsDictionary)
	c.Assert(result.ok, Equals, true)

	c.Assert(len(spec.scenarios[0].steps), Equals, 2)

	firstConcept := spec.scenarios[0].steps[0]
	c.Assert(firstConcept.isConcept, Equals, true)
	c.Assert(firstConcept.conceptSteps[0].value, Equals, "enter user {} and {}")
	c.Assert(firstConcept.conceptSteps[0].args[0].argType, Equals, dynamic)
	c.Assert(firstConcept.conceptSteps[0].args[0].value, Equals, "user-id")
	c.Assert(firstConcept.conceptSteps[0].args[1].argType, Equals, dynamic)
	c.Assert(firstConcept.conceptSteps[0].args[1].value, Equals, "user-description")
	c.Assert(firstConcept.conceptSteps[1].value, Equals, "select {}")
	c.Assert(firstConcept.conceptSteps[1].args[0].value, Equals, "finish")
	c.Assert(firstConcept.conceptSteps[1].args[0].argType, Equals, static)

	c.Assert(len(firstConcept.lookup.paramValue), Equals, 2)
	c.Assert(firstConcept.lookup.paramValue[0].name, Equals, "user-id")
	c.Assert(firstConcept.lookup.paramValue[0].value, Equals, "id")
	c.Assert(firstConcept.lookup.paramValue[0].paramType, Equals, dynamic)
	c.Assert(firstConcept.lookup.paramValue[1].name, Equals, "user-description")
	c.Assert(firstConcept.lookup.paramValue[1].value, Equals, "description")
	c.Assert(firstConcept.lookup.paramValue[1].paramType, Equals, dynamic)

	secondConcept := spec.scenarios[0].steps[1]
	c.Assert(secondConcept.isConcept, Equals, true)
	c.Assert(secondConcept.conceptSteps[0].value, Equals, "enter user {} and {}")
	c.Assert(secondConcept.conceptSteps[0].args[0].argType, Equals, dynamic)
	c.Assert(secondConcept.conceptSteps[0].args[0].value, Equals, "user-id")
	c.Assert(secondConcept.conceptSteps[0].args[1].argType, Equals, dynamic)
	c.Assert(secondConcept.conceptSteps[0].args[1].value, Equals, "user-description")
	c.Assert(secondConcept.conceptSteps[1].value, Equals, "select {}")
	c.Assert(secondConcept.conceptSteps[1].args[0].value, Equals, "finish")
	c.Assert(secondConcept.conceptSteps[1].args[0].argType, Equals, static)

	c.Assert(len(secondConcept.lookup.paramValue), Equals, 2)
	c.Assert(secondConcept.lookup.paramValue[0].name, Equals, "user-id")
	c.Assert(secondConcept.lookup.paramValue[0].value, Equals, "456")
	c.Assert(secondConcept.lookup.paramValue[0].paramType, Equals, static)
	c.Assert(secondConcept.lookup.paramValue[1].name, Equals, "user-description")
	c.Assert(secondConcept.lookup.paramValue[1].value, Equals, "Regular fellow")
	c.Assert(secondConcept.lookup.paramValue[1].paramType, Equals, static)

}
