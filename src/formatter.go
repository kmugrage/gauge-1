package main

import (
	"bytes"
	"fmt"
	//"math"
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

func formatStepText(stepText string) string {
	return fmt.Sprintf("* %s", stepText)
}

func formatHeading(heading, headingChar string) string {
	length := len(heading)
	if length > 10 {
		length = 10
	}

	return fmt.Sprintf("%s\n%s\n", heading, getRepeatedChars(headingChar, length))
}

func formatTable(table *table) string {
	columnToWidthMap := make(map[int]int)
	for i, header := range table.headers {
		//table.get(header) returns a list of cells in that particular column
		cells := table.get(header)
		columnToWidthMap[i] = findLongestCellWidth(cells, len(header))
	}

	var tableStringBuffer bytes.Buffer
	tableStringBuffer.WriteString("|")
	for i, header := range table.headers {
		width := columnToWidthMap[i]
		tableStringBuffer.WriteString(fmt.Sprintf("%s|", addPaddingToCell(header, width)))
	}

	tableStringBuffer.WriteString("\n")
	tableStringBuffer.WriteString("|")
	for i, _ := range table.headers {
		width := columnToWidthMap[i]
		cell := getRepeatedChars("-", width)
		tableStringBuffer.WriteString(fmt.Sprintf("%s|", addPaddingToCell(cell, width)))
	}

	tableStringBuffer.WriteString("\n")
	for _, row := range table.getRows() {
		tableStringBuffer.WriteString("|")
		for i, cell := range row {
			width := columnToWidthMap[i]
			tableStringBuffer.WriteString(fmt.Sprintf("%s|", addPaddingToCell(cell, width)))
		}
		tableStringBuffer.WriteString("\n")
	}

	return string(tableStringBuffer.Bytes())
}

func addPaddingToCell(cellValue string, width int) string {
	padding := getRepeatedChars(" ", width-len(cellValue))
	return fmt.Sprintf("%s%s", cellValue, padding)
}

func findLongestCellWidth(columnCells []tableCell, minValue int) int {
	longestLength := minValue
	for _, cellValue := range columnCells {
		cellValueLen := len(cellValue.value)
		if cellValueLen > longestLength {
			longestLength = cellValueLen
		}
	}
	return longestLength
}
