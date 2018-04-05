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
		if len(l) > 80 {
			l = l[0:80]
		}

		report[index] = Line{
			Text: []rune(fmt.Sprintf("%s", l)),
		}
	}
	return report
}
