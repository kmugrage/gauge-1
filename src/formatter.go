package main

import (
	"bytes"
	"fmt"
	"math"
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
	columnIndexToWidthMap := make(map[int]int)
	for i, header := range table.headers {
		//table.get(header) returns a list of cells in that particular column
		columnCells := table.get(header)
		columnIndexToWidthMap[i] = longestColumnCellLength(columnCells, len(header))
	}

	var tableStringBuffer *bytes.Buffer
    tableStringBuffer.WriteString("|")
    for i, header := range table.headers {
        width := columnIndexToWidthMap[i]
        centerToCell(header, width)
        tableStringBuffer.WriteString(fmt.Sprintf("%s|"))
    }
    for _, table.getRowCount()


}


func centerToCell(cellValue string, width int) {
    paddinglen := math.Floor((len(cellValue) - width)/2)
    padding := getRepeatedChars(" ",int(paddinglen))
    padded := fmt.Sprintf("%s%s%s", padding, cellValue, padding))
    // When paddingLen is odd
    if(len(padded) != width) {
        padded =  fmt.Sprintf("%s ",padded )
    } 
    return padded
}

func longestColumnCellLength(columnCells []string, minValue int) int {
	longestLength := minValue
	for _, cellValue := range columnCells {
		cellValueLen = len(cellValue)
		if cellValueLen > longestLength {
			longestLength = cellValueLen
		}
	}
	return longestLength
}
