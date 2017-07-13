package main

import (
	"github.com/edi9999/ii/core"
	"github.com/edi9999/ii/interactive"
	"github.com/gdamore/tcell"
)

type VisualState struct {
	folders        []interactive.Line
	selected       int
	xbound, ybound int
}

func NewVisualState(state core.State) VisualState {
	lines := interactive.PrintTree(state.Buffers[0].Lines)
	xbound := 0
	ybound := len(lines)
	for index, line := range lines {
		if len(line.Text)-1 > xbound {
			xbound = len(line.Text) - 1
		}
		lines[index] = line
	}
	return VisualState{lines, 0, xbound, ybound}
}

func (vs VisualState) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	style := tcell.StyleDefault
	if y == vs.selected {
		style = style.Reverse(true)
	}
	if y < len(vs.folders) {
		line := vs.folders[y]
		if x < len(vs.folders[y].Text) {
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
