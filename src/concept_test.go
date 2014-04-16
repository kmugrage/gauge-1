package main

import (
	. "launchpad.net/gocheck"
)

func (s *MySuite) TestParsingSimpleConcept(c *C) {
	parser := new(conceptParser)
	concepts, err := parser.parse("# my concept \n * first step \n * second step ")

	c.Assert(err, Equals, nil)
	c.Assert(len(concepts), Equals, 1)

	concept := concepts[0]

	c.Assert(concept.isConcept, Equals, true)
	c.Assert(len(concept.conceptSteps), Equals, 2)
	c.Assert(concept.conceptSteps[0].value, Equals, "first step")
	c.Assert(concept.conceptSteps[1].value, Equals, "second step")

}

func (s *MySuite) TestParsingSimpleConceptWithParamters(c *C) {
	parser := new(conceptParser)
	concepts, err := parser.parse("# my concept with <param0> and <param1> \n * first step using <param0> \n * second step using \"value\" and <param1> ")

	c.Assert(err, Equals, nil)
	c.Assert(len(concepts), Equals, 1)

	concept := concepts[0]
	c.Assert(concept.isConcept, Equals, true)
	c.Assert(len(concept.conceptSteps), Equals, 2)
	c.Assert(len(concept.lookup.paramValue), Equals, 2)
	c.Assert(concept.lookup.containsParam("param0"), Equals, true)
	c.Assert(concept.lookup.containsParam("param1"), Equals, true)

	firstStep := concept.conceptSteps[0]
	c.Assert(firstStep.value, Equals, "first step using {}")
	c.Assert(len(firstStep.args), Equals, 1)
	c.Assert(firstStep.args[0].argType, Equals, dynamic)
	c.Assert(firstStep.args[0].value, Equals, "param0")

	secondStep := concept.conceptSteps[1]
	c.Assert(secondStep.value, Equals, "second step using {} and {}")
	c.Assert(len(secondStep.args), Equals, 2)
	c.Assert(secondStep.args[0].argType, Equals, static)
	c.Assert(secondStep.args[0].value, Equals, "value")
	c.Assert(secondStep.args[1].argType, Equals, dynamic)
	c.Assert(secondStep.args[1].value, Equals, "param1")

}
