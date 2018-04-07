package core

import (
	"fmt"
	"github.com/gdamore/tcell"
)

type VisualState struct {
	lines          []Line
	style          tcell.Style
	selected       int
	xbound, ybound int
}

type Line struct {
	Text []rune
}

func printLines(lines []string) []Line {
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

func NewVisualState(linesArray []string, style tcell.Style) VisualState {
	lines := printLines(linesArray)
	xbound := 0
	ybound := len(lines)
	for index, line := range lines {
		if len(line.Text)-1 > xbound {
			xbound = len(line.Text) - 1
		}
		lines[index] = line
	}
	return VisualState{lines, style, 0, xbound, ybound}
}

func (vs VisualState) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	style := vs.style
	if y < len(vs.lines) {
		line := vs.lines[y]
		if x < len(vs.lines[y].Text) {
			return line.Text[x], style, nil, 1
		}
	}
	return ' ', style, nil, 1
}
func (vs VisualState) GetBounds() (int, int) {
	return vs.xbound, vs.ybound
}
func (VisualState) SetCursor(int, int) {
}

func (VisualState) GetCursor() (int, int, bool, bool) {
	return 0, 0, false, false
}
func (VisualState) MoveCursor(offx, offy int) {
}
