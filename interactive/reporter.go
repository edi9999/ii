package interactive

import (
	"fmt"
)

type Line struct {
	Text []rune
}

func PrintTree(lines []string) []Line {
	report := make([]Line, len(lines))
	for index, l := range lines {
		report[index] = Line{
			Text: []rune(fmt.Sprintf("%s", l)),
		}
	}
	return report
}
