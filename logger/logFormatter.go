// Copyright 2015 ThoughtWorks, Inc.

// This file is part of Gauge.

// Gauge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Gauge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Gauge.  If not, see <http://www.gnu.org/licenses/>.

package logger

import (
	"fmt"
	"strings"
)

const (
	scenarioIndentation = 2
	stepIndentation     = 4
	sysoutIndentation   = 8
)

func formatScenario(msg string) string {
	return fmt.Sprintf("## %s", msg)
}

func formatStep(msg string) string {
	return fmt.Sprintf("%s", msg)
}

func formatSpec(msg string) string {
	return fmt.Sprintf("# %s", msg)
}

func indent(text string, indentation int) string {
	return spaces(indentation) + strings.TrimSpace(text)
}

func spaces(numOfSpaces int) string {
	text := ""
	for i := 0; i < numOfSpaces; i++ {
		text += " "
	}
	return text
}
