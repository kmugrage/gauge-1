package logger

import (
	"bytes"
	"fmt"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/op/go-logging"
	"strings"
)

const (
	success = "✔ "
	failure = "✘ "
)

type coloredLogger struct {
	heading   string
	errBuffer bytes.Buffer
}

func newColoredConsoleWriter() *coloredLogger {
	return &coloredLogger{}
}

func (cLogger *coloredLogger) writeSysoutBuffer(text string) {
	if level == logging.DEBUG {
		text = strings.Replace(text, "\n", "\n\t", -1)
		cLogger.Write(stepIndentation, text, ct.None, false)
	}
}

func (cLogger *coloredLogger) SpecStart(heading string) {
	msg := formatSpec(heading)
	Log.Info(msg)
	ct.Foreground(ct.Cyan, true)
	ConsoleWrite(msg)
	fmt.Println()
	ct.ResetColor()
}

func (coloredLogger *coloredLogger) SpecEnd() {
}

func (cLogger *coloredLogger) ScenarioStart(scenarioHeading string) {
	msg := formatScenario(scenarioHeading)
	Log.Info(msg)

	cLogger.heading = strings.Trim(msg, "\n")
	if level == logging.INFO {
		cLogger.Print(scenarioIndentation, cLogger.heading+spaces(4), ct.Yellow, false)
	} else {
		cLogger.Write(scenarioIndentation, cLogger.heading, ct.Magenta, false)
		fmt.Println()
	}
}

func (cLogger *coloredLogger) ScenarioEnd(failed bool) {
	if level == logging.INFO {
		cLogger.Write(stepIndentation, "\n"+cLogger.errBuffer.String(), ct.Red, false)
	}
	cLogger.heading = ""
	cLogger.errBuffer.Reset()
}

func (cLogger *coloredLogger) StepStart(stepText string) {
	Log.Debug(stepText)
	if level == logging.DEBUG {
		cLogger.heading = strings.Trim(stepText, "\n")
		cLogger.Write(stepIndentation, "Executing => "+cLogger.heading, ct.Yellow, false)
	}
}

func (cLogger *coloredLogger) StepEnd(failed bool) {
	if level == logging.DEBUG {
		if failed {
			cLogger.Write(stepIndentation, "[Fail]", ct.Red, false)
		} else {
			cLogger.Write(stepIndentation, "[Pass]", ct.Green, false)
		}
		fmt.Println()
	} else {
		if failed {
			cLogger.Print(0, failure, ct.Red, false)
		} else {
			cLogger.Print(0, success, ct.Green, false)
		}
	}
}

func (cLogger *coloredLogger) Error(msg string, args ...interface{}) {
	Log.Error(msg, args)
	if level == logging.INFO {
		cLogger.errBuffer.WriteString(fmt.Sprintf(msg+"\n", args))
	} else {
		cLogger.Write(stepIndentation, fmt.Sprintf(msg, args), ct.Red, false)
	}
}

func (cLogger *coloredLogger) Write(indentation int, text string, color ct.Color, isBright bool) {
	ct.Foreground(color, isBright)
	fmt.Println(spaces(indentation) + text)
	ct.ResetColor()
}

func (cLogger *coloredLogger) Print(indentation int, text string, color ct.Color, isBright bool) {
	ct.Foreground(color, isBright)
	fmt.Print(spaces(indentation) + text)
	ct.ResetColor()
}
