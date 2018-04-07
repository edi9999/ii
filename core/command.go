package core

import (
	"github.com/edi9999/ii/events"
	"regexp"
	"strings"
)

type action struct {
	t actionType
	a string
}
type actionType int

const (
	actIgnore actionType = iota
	actInvalid
	actRune
	actMouse
	actBeginningOfLine
	actAbort
	actAccept
	actBackwardChar
	actBackwardDeleteChar
	actBackwardWord
	actCancel
	actClearScreen
	actDeleteChar
	actDeleteCharEOF
	actEndOfLine
	actForwardChar
	actForwardWord
	actKillLine
	actKillWord
	actUnixLineDiscard
	actUnixWordRubout
	actYank
	actBackwardKillWord
	actSelectAll
	actDeselectAll
	actToggle
	actToggleAll
	actToggleDown
	actToggleUp
	actToggleIn
	actToggleOut
	actDown
	actUp
	actPageUp
	actPageDown
	actHalfPageUp
	actHalfPageDown
	actJump
	actJumpAccept
	actPrintQuery
	actToggleSort
	actTogglePreview
	actTogglePreviewWrap
	actPreviewUp
	actPreviewDown
	actPreviewPageUp
	actPreviewPageDown
	actPreviousHistory
	actNextHistory
	actExecute
	actExecuteSilent
	actExecuteMulti // Deprecated
	actSigStop
	actTop
)

func (t *LineInput) delChar() bool {
	if len(t.Input) > 0 && t.Cx < len(t.Input) {
		t.Input = append(t.Input[:t.Cx], t.Input[t.Cx+1:]...)
		return true
	}
	return false
}

func toActions(types ...actionType) []action {
	actions := make([]action, len(types))
	for idx, t := range types {
		actions[idx] = action{t: t, a: ""}
	}
	return actions
}
func copySlice(slice []rune) []rune {
	ret := make([]rune, len(slice))
	copy(ret, slice)
	return ret
}
func findLastMatch(pattern string, str string) int {
	rx, err := regexp.Compile(pattern)
	if err != nil {
		return -1
	}
	locs := rx.FindAllStringIndex(str, -1)
	if locs == nil {
		return -1
	}
	return locs[len(locs)-1][0]
}
func findFirstMatch(pattern string, str string) int {
	rx, err := regexp.Compile(pattern)
	if err != nil {
		return -1
	}
	loc := rx.FindStringIndex(str)
	if loc == nil {
		return -1
	}
	return loc[0]
}

func (t *LineInput) rubout(pattern string) {
	pCx := t.Cx
	after := t.Input[t.Cx:]
	t.Cx = findLastMatch(pattern, string(t.Input[:t.Cx])) + 1
	t.Yanked = copySlice(t.Input[t.Cx:pCx])
	t.Input = append(t.Input[:t.Cx], after...)
}

func keyMatch(key int, event Event) bool {
	return event.Type == key ||
		event.Type == events.Rune && int(event.Char) == key-events.AltZ ||
		event.Type == events.Mouse && key == events.DoubleClick && event.MouseEvent.Double
}

func quoteEntry(entry string) string {
	return "'" + strings.Replace(entry, "'", "'\\''", -1) + "'"
}

func defaultKeymap() map[int][]action {
	keymap := make(map[int][]action)
	keymap[events.Invalid] = toActions(actInvalid)
	keymap[events.Resize] = toActions(actClearScreen)
	keymap[events.CtrlA] = toActions(actBeginningOfLine)
	keymap[events.CtrlB] = toActions(actBackwardChar)
	keymap[events.CtrlC] = toActions(actAbort)
	keymap[events.CtrlG] = toActions(actAbort)
	keymap[events.CtrlQ] = toActions(actAbort)
	keymap[events.ESC] = toActions(actAbort)
	keymap[events.CtrlD] = toActions(actDeleteCharEOF)
	keymap[events.CtrlE] = toActions(actEndOfLine)
	keymap[events.CtrlF] = toActions(actForwardChar)
	keymap[events.CtrlH] = toActions(actBackwardDeleteChar)
	keymap[events.BSpace] = toActions(actBackwardDeleteChar)
	keymap[events.Tab] = toActions(actToggleDown)
	keymap[events.BTab] = toActions(actToggleUp)
	keymap[events.CtrlJ] = toActions(actDown)
	keymap[events.CtrlK] = toActions(actUp)
	keymap[events.CtrlL] = toActions(actClearScreen)
	keymap[events.CtrlM] = toActions(actAccept)
	keymap[events.CtrlN] = toActions(actDown)
	keymap[events.CtrlP] = toActions(actUp)
	keymap[events.CtrlU] = toActions(actUnixLineDiscard)
	keymap[events.CtrlW] = toActions(actUnixWordRubout)
	keymap[events.CtrlY] = toActions(actYank)

	keymap[events.AltB] = toActions(actBackwardWord)
	keymap[events.SLeft] = toActions(actBackwardWord)
	keymap[events.AltF] = toActions(actForwardWord)
	keymap[events.SRight] = toActions(actForwardWord)
	keymap[events.AltD] = toActions(actKillWord)
	keymap[events.AltBS] = toActions(actBackwardKillWord)

	keymap[events.Up] = toActions(actUp)
	keymap[events.Down] = toActions(actDown)
	keymap[events.Left] = toActions(actBackwardChar)
	keymap[events.Right] = toActions(actForwardChar)

	keymap[events.Home] = toActions(actBeginningOfLine)
	keymap[events.End] = toActions(actEndOfLine)
	keymap[events.Del] = toActions(actDeleteChar)
	keymap[events.PgUp] = toActions(actPageUp)
	keymap[events.PgDn] = toActions(actPageDown)

	keymap[events.Rune] = toActions(actRune)
	keymap[events.Mouse] = toActions(actMouse)
	keymap[events.DoubleClick] = toActions(actAccept)
	return keymap
}

type Buf struct {
	Lines  []string
	Status int
	Cmd    string
	Index  int
}

type State struct {
	Buffers        []Buf
	ExitCodes      []int
	LineInput      LineInput
	Stdin          []string
	SelectedWidget int
}

type LineInput struct {
	Input  []rune
	Cx     int
	Yanked []rune
}

// Executer represents a user action triggered on a State
type Executer interface {
	Execute(State) (State, error)
}

type Event struct {
	Type       int
	Char       rune
	MouseEvent *events.MouseEvent
}

func copyState(state State) State {
	return State{
		Buffers:   state.Buffers,
		LineInput: state.LineInput,
		Stdin:     state.Stdin,
	}
}

func (e Event) Execute(oldState State) (State, error) {
	keymap := defaultKeymap()
	mapkey := e.Type
	actions := keymap[e.Type]
	newState := copyState(oldState)
	t := newState.LineInput

	var doAction func(action, int) bool
	doActions := func(actions []action, mapkey int) bool {
		for _, action := range actions {
			if !doAction(action, mapkey) {
				return false
			}
		}
		return true
	}
	wordRubout := "[^[:alnum:]][[:alnum:]]"
	wordNext := "[[:alnum:]][^[:alnum:]]|(.$)"
	doAction = func(a action, mapkey int) bool {
		switch a.t {
		case actIgnore:
		case actInvalid:
			return false
		case actBeginningOfLine:
			t.Cx = 0
		case actBackwardChar:
			if t.Cx > 0 {
				t.Cx--
			}
		case actDeleteChar:
			t.delChar()
		case actEndOfLine:
			t.Cx = len(t.Input)
		case actCancel:
			if len(t.Input) != 0 {
				t.Yanked = t.Input
				t.Input = []rune{}
				t.Cx = 0
			}
		case actForwardChar:
			if t.Cx < len(t.Input) {
				t.Cx++
			}
		case actBackwardDeleteChar:
			if t.Cx > 0 {
				t.Input = append(t.Input[:t.Cx-1], t.Input[t.Cx:]...)
				t.Cx--
			}
		case actUnixLineDiscard:
			if t.Cx > 0 {
				t.Yanked = copySlice(t.Input[:t.Cx])
				t.Input = t.Input[t.Cx:]
				t.Cx = 0
			}
		case actUnixWordRubout:
			if t.Cx > 0 {
				t.rubout("\\s\\S")
			}
		case actBackwardKillWord:
			if t.Cx > 0 {
				t.rubout(wordRubout)
			}
		case actYank:
			suffix := copySlice(t.Input[t.Cx:])
			t.Input = append(append(t.Input[:t.Cx], t.Yanked...), suffix...)
			t.Cx += len(t.Yanked)
		case actBackwardWord:
			t.Cx = findLastMatch(wordRubout, string(t.Input[:t.Cx])) + 1
		case actForwardWord:
			t.Cx += findFirstMatch(wordNext, string(t.Input[t.Cx:])) + 1
		case actKillWord:
			nCx := t.Cx +
				findFirstMatch(wordNext, string(t.Input[t.Cx:])) + 1
			if nCx > t.Cx {
				t.Yanked = copySlice(t.Input[t.Cx:nCx])
				t.Input = append(t.Input[:t.Cx], t.Input[nCx:]...)
			}
		case actKillLine:
			if t.Cx < len(t.Input) {
				t.Yanked = copySlice(t.Input[t.Cx:])
				t.Input = t.Input[:t.Cx]
			}
		case actRune:
			prefix := copySlice(t.Input[:t.Cx])
			t.Input = append(append(prefix, e.Char), t.Input[t.Cx:]...)
			t.Cx++
		}
		return true
	}
	doActions(actions, mapkey)
	newState.LineInput = t

	return newState, nil
}
