package main

import (
	"github.com/edi9999/ii/core"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"strconv"
	"strings"
	"sync"
)

func InteractiveTree(s tcell.Screen, states chan core.State, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		state, more := <-states
		if !more {
			break
		}
		printOptions(state, s)
	}
}

func printOptions(state core.State, s tcell.Screen) {
	s.Clear()
	_, height := s.Size()
	layout := views.NewBoxLayout(views.Vertical)
	statusBar := views.NewSimpleStyledText()
	statusBar.RegisterStyle('n', tcell.StyleDefault)
	statusBar.RegisterStyle('r', tcell.StyleDefault.
		Foreground(tcell.ColorRed))
	statusBar.RegisterStyle('g', tcell.StyleDefault.
		Foreground(tcell.ColorGreen))

	mystr := ""
	subLayout := views.NewBoxLayout(views.Horizontal)
	widgets := []*views.CellView{}

	color := ""
	if len(state.Stdin) > 0 {
		middle := views.NewCellView()
		middle.SetModel(NewVisualState(state.Stdin))
		widgets = append(widgets, middle)
	}
	for i, buffer := range state.Buffers {
		color = "%r"
		switch buffer.Status {
		case 0:
			color = "%g"
		case 1:
			color = "%r"
		}
		mystr = mystr + color + strconv.Itoa(buffer.Status) + strings.Repeat(" ", len(buffer.Cmd))
		middle := views.NewCellView()
		middle.SetModel(NewVisualState(state.Buffers[i].Lines))
		widgets = append(widgets, middle)
	}
	statusBar.SetMarkup(mystr)

	text := views.NewTextArea()
	text.SetContent(string(state.LineInput.Input))

	layout.SetView(s)

	for _, widget := range widgets {
		subLayout.AddWidget(widget, 1)
	}
	for _, widget := range widgets {
		widget.Resize()
	}

	layout.AddWidget(subLayout, 1)
	subLayout.Resize()
	layout.AddWidget(statusBar, 0.0)
	layout.AddWidget(text, 0.0)

	text.Resize()

	layout.Draw()
	s.ShowCursor(state.LineInput.Cx, height-1)
	s.Show()
}
