package main

import (
	"bytes"
	"fmt"
	"github.com/wsxiaoys/terminal"
	"os"
	"strings"
)

const (
	mode_unbuffered = 0
	mode_buffered   = 1
)

type consoleWriter struct {
	mode   int
	buffer *bytes.Buffer
}

func newConsoleWriter() *consoleWriter {
	var b bytes.Buffer
	return &consoleWriter{buffer: &b, mode: mode_unbuffered}
}

var currentConsoleWriter *consoleWriter

func getCurrentConsole() *consoleWriter {
	if currentConsoleWriter == nil {
		currentConsoleWriter = newConsoleWriter()
	}
	return currentConsoleWriter
}

func (writer *consoleWriter) enableBuffering() {
	writer.mode = mode_buffered
}

func (writer *consoleWriter) disableBuffering() {
	writer.mode = mode_unbuffered
}

func (writer *consoleWriter) Write(b []byte) (int, error) {
	length := 0
	if writer.mode == mode_unbuffered {
		length = writer.writeToStdout(b)
	} else {
		_, err := writer.buffer.Write(b)
		if err != nil {
			writer.writeToStdout([]byte(fmt.Sprintf("[Error] Failed to buffer. %s\n", err.Error())))
			writer.writeToStdout(b)
		}
	}

	return length, nil
}

func (writer *consoleWriter) writeToStdout(b []byte) int {
	length, err := os.Stdout.Write(b)
	if err != nil {
		panic(err)
	}

	return length
}

func (writer *consoleWriter) flush() (int, error) {
	i, err := os.Stdout.Write(writer.buffer.Bytes())
	if err != nil {
		writer.buffer.Reset()
	}
	return i, err
}

func (writer *consoleWriter) writeSpecExecutionHeading(specHeading string) {
	formattedHeading := formatSpecHeading(specHeading)
	writer.Write([]byte(formattedHeading))
}

func (writer *consoleWriter) writeScenarioHeading(scenarioHeading string) {
	formattedHeading := formatScenarioHeading(scenarioHeading)
	writer.Write([]byte(formattedHeading))
}

func (writer *consoleWriter) writeStep(stepRequest *ExecuteStepRequest) {
	terminal.Stdout.Colorf("@b%s\n", formatStepRequest(stepRequest))
	writer.enableBuffering()
}

func extractStepWithResolvedParameters(stepRequest *ExecuteStepRequest) string {
	text := stepRequest.GetParsedStepText()
	paramCount := strings.Count(text, PARAMETER_PLACEHOLDER)
	for i := 0; i < paramCount; i++ {
		text = strings.Replace(text, PARAMETER_PLACEHOLDER, resolveParameterText(stepRequest.GetArgs()[i]), 1)
	}
	return text
}

func resolveParameterText(argument *Argument) string {
	if argument.GetType() == "table" {
		table := tableFrom(argument.GetTable())
		formattedTable := formatTable(table)
		return fmt.Sprintf("\n%s", formattedTable)
	}
	return fmt.Sprintf("\"%s\"", argument.GetValue())
}

func (writer *consoleWriter) writeStepFinished(stepRequest *ExecuteStepRequest, isPassed bool) {
	stepText := formatStepRequest(stepRequest)
	terminal.Stdout.Up(strings.Count(stepText, "\n") + 1)
	if isPassed {
		terminal.Stdout.Colorf("@g%s\n", stepText)
	} else {
		terminal.Stdout.Colorf("@r%s\n", stepText)
	}
	writer.flush()
	writer.disableBuffering()
}

func formatStepRequest(stepRequest *ExecuteStepRequest) string {
	return formatStepText(extractStepWithResolvedParameters(stepRequest))
}
