package tui

import (
	"encoding/json"
	"fmt"
	"github.com/edi9999/ii/core"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"io/ioutil"
	"log"
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

var baseString = strings.Split(strings.Repeat("\n", 100), "\n")

func createBorderWidget() *views.CellView {
	border := views.NewCellView()
	border.SetModel(core.NewVisualState(baseString, tcell.StyleDefault.Background(tcell.NewRGBColor(44, 44, 44))))
	return border
}

func createSelectedBorderWidget() *views.CellView {
	border := views.NewCellView()
	border.SetModel(core.NewVisualState(baseString, tcell.StyleDefault.Background(tcell.Color32)))
	return border
}

func printOptions(state core.State, s tcell.Screen) {
	s.Clear()
	intWidth, height := s.Size()
	width := float64(intWidth)

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
	countViews := len(state.Buffers)
	selectedWidget := state.SelectedWidget

	color := ""
	if len(state.Stdin) > 0 {
		countViews = countViews + 1
		selectedWidget = selectedWidget + 1
		middle := views.NewCellView()
		middle.SetModel(core.NewVisualState(state.Stdin, tcell.StyleDefault))
		widgets = append(widgets, middle)
	}
	lastIndex := 0
	for i, buffer := range state.Buffers {
		color = "%r"
		switch buffer.Status {
		case 0:
			color = "%g"
		case 1:
			color = "%r"
		}
		mystr = mystr + strings.Repeat(" ", buffer.Index-lastIndex) + color + strconv.Itoa(buffer.Status)
		lastIndex = buffer.Index
		middle := views.NewCellView()
		middle.SetModel(core.NewVisualState(state.Buffers[i].Lines, tcell.StyleDefault))
		widgets = append(widgets, middle)
	}
	statusBar.SetMarkup(mystr)

	text := views.NewTextArea()
	text.SetContent(string(state.LineInput.Input))

	layout.SetView(s)

	borderWidgets := []*views.CellView{}

	if countViews == 0 {
		countViews = 1
	}

	if selectedWidget > countViews-1 {
		selectedWidget = countViews - 1
	}
	if state.LineInput.Cx == 0 {
		selectedWidget = 0
	}
	selectedRatio := 1.5
	if len(widgets) == 1 {
		selectedRatio = 1
	}
	selectedWidth := width/selectedRatio - 1
	widgetWidth := (width*(1-1/selectedRatio))/float64(countViews) - 1
	for i, widget := range widgets {
		w := widgetWidth
		var border *views.CellView
		if selectedWidget == i {
			border = createSelectedBorderWidget()
			w = selectedWidth
		} else {
			border = createBorderWidget()
		}
		borderWidgets = append(borderWidgets, border)

		subLayout.AddWidget(border, 1)
		subLayout.AddWidget(widget, float64(w))
	}
	for _, widget := range subLayout.Widgets() {
		widget.Resize()
	}

	layout.AddWidget(subLayout, 1)
	subLayout.Resize()
	layout.AddWidget(statusBar, 0.0)
	layout.AddWidget(text, 0.0)

	text.Resize()

	data, err := json.Marshal(state)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("/tmp/ii.log", []byte(fmt.Sprintf("%s\n", data)), 0644)

	layout.Draw()
	s.ShowCursor(state.LineInput.Cx, height-1)
	s.Show()
}
