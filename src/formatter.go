package main

import (
	"fmt"
)

func getRepeatedChars(character string, repeatCount int) string {
	formatted := ""
	for i := 0; i < repeatCount; i++ {
		formatted = fmt.Sprintf("%s%s", formatted, character)
	}

	return formatted
}

func formatSpecHeading(specHeading string) string {
	return formatHeading(specHeading, "=")
}

func formatScenarioHeading(scenarioHeading string) string {
	return formatHeading(scenarioHeading, "-")
}

func formatHeading(heading, headingChar string) string {
	length := len(heading)
	if length > 10 {
		length = 10
	}

	return fmt.Sprintf("%s\n%s\n", heading, getRepeatedChars(headingChar, length))
}
