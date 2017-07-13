package main

import (
	"github.com/edi9999/ii/core"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
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
	inner := views.NewBoxLayout(views.Vertical)
	middle := views.NewCellView()
	text := views.NewTextArea()
	text.SetContent(string(state.LineInput.Input))
	s.ShowCursor(state.LineInput.Cx, 53)
	middle.SetModel(NewVisualState(state))
	inner.SetView(s)
	inner.AddWidget(middle, 0.33)
	inner.AddWidget(text, 0.0)
	middle.Resize()
	text.Resize()
	inner.Draw()
	s.Show()
}
