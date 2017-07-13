package core

import (
	"github.com/edi9999/ii/tui"
	"log"
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
		event.Type == tui.Rune && int(event.Char) == key-tui.AltZ ||
		event.Type == tui.Mouse && key == tui.DoubleClick && event.MouseEvent.Double
}

func quoteEntry(entry string) string {
	return "'" + strings.Replace(entry, "'", "'\\''", -1) + "'"
}

func defaultKeymap() map[int][]action {
	keymap := make(map[int][]action)
	keymap[tui.Invalid] = toActions(actInvalid)
	keymap[tui.Resize] = toActions(actClearScreen)
	keymap[tui.CtrlA] = toActions(actBeginningOfLine)
	keymap[tui.CtrlB] = toActions(actBackwardChar)
	keymap[tui.CtrlC] = toActions(actAbort)
	keymap[tui.CtrlG] = toActions(actAbort)
	keymap[tui.CtrlQ] = toActions(actAbort)
	keymap[tui.ESC] = toActions(actAbort)
	keymap[tui.CtrlD] = toActions(actDeleteCharEOF)
	keymap[tui.CtrlE] = toActions(actEndOfLine)
	keymap[tui.CtrlF] = toActions(actForwardChar)
	keymap[tui.CtrlH] = toActions(actBackwardDeleteChar)
	keymap[tui.BSpace] = toActions(actBackwardDeleteChar)
	keymap[tui.Tab] = toActions(actToggleDown)
	keymap[tui.BTab] = toActions(actToggleUp)
	keymap[tui.CtrlJ] = toActions(actDown)
	keymap[tui.CtrlK] = toActions(actUp)
	keymap[tui.CtrlL] = toActions(actClearScreen)
	keymap[tui.CtrlM] = toActions(actAccept)
	keymap[tui.CtrlN] = toActions(actDown)
	keymap[tui.CtrlP] = toActions(actUp)
	keymap[tui.CtrlU] = toActions(actUnixLineDiscard)
	keymap[tui.CtrlW] = toActions(actUnixWordRubout)
	keymap[tui.CtrlY] = toActions(actYank)

	keymap[tui.AltB] = toActions(actBackwardWord)
	keymap[tui.SLeft] = toActions(actBackwardWord)
	keymap[tui.AltF] = toActions(actForwardWord)
	keymap[tui.SRight] = toActions(actForwardWord)
	keymap[tui.AltD] = toActions(actKillWord)
	keymap[tui.AltBS] = toActions(actBackwardKillWord)

	keymap[tui.Up] = toActions(actUp)
	keymap[tui.Down] = toActions(actDown)
	keymap[tui.Left] = toActions(actBackwardChar)
	keymap[tui.Right] = toActions(actForwardChar)

	keymap[tui.Home] = toActions(actBeginningOfLine)
	keymap[tui.End] = toActions(actEndOfLine)
	keymap[tui.Del] = toActions(actDeleteChar)
	keymap[tui.PgUp] = toActions(actPageUp)
	keymap[tui.PgDn] = toActions(actPageDown)

	keymap[tui.Rune] = toActions(actRune)
	keymap[tui.Mouse] = toActions(actMouse)
	keymap[tui.DoubleClick] = toActions(actAccept)
	return keymap
}

type Buf struct {
	Lines []string
}

type State struct {
	Buffers   []Buf
	LineInput LineInput
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
	MouseEvent *tui.MouseEvent
}

func copyState(state State) State {
	return State{
		Buffers:   state.Buffers,
		LineInput: state.LineInput,
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
			log.Printf("%v\n", actions)
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
