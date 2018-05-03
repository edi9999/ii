package events

import (
	"github.com/gdamore/tcell"
	"strconv"
)

type Attr tcell.Style

// Types of user action
const (
	Rune = iota

	CtrlA
	CtrlB
	CtrlC
	CtrlD
	CtrlE
	CtrlF
	CtrlG
	CtrlH
	Tab
	CtrlJ
	CtrlK
	CtrlL
	CtrlM
	CtrlN
	CtrlO
	CtrlP
	CtrlQ
	CtrlR
	CtrlS
	CtrlT
	CtrlU
	CtrlV
	CtrlW
	CtrlX
	CtrlY
	CtrlZ
	ESC
	CtrlSpace

	Invalid
	Resize
	Mouse
	DoubleClick

	BTab
	BSpace

	Del
	PgUp
	PgDn

	Up
	Down
	Left
	Right
	Home
	End

	SLeft
	SRight
	SUp
	SDown

	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12

	Change

	AltSpace
	AltSlash
	AltBS

	Alt0
)

const ( // Reset iota
	AltA = Alt0 + 'a' - '0' + iota
	AltB
	AltC
	AltD
	AltE
	AltF
	AltZ     = AltA + 'z' - 'a'
	CtrlAltA = AltZ + 1
	CtrlAltM = CtrlAltA + 'm' - 'a'
)

type Color int32

func (c Color) is24() bool {
	return c > 0 && (c&(1<<24)) > 0
}

const (
	colUndefined Color = -2
	colDefault         = -1
)

const (
	colBlack Color = iota
	colRed
	colGreen
	colYellow
	colBlue
	colMagenta
	colCyan
	colWhite
)

type FillReturn int

type ColorPair struct {
	fg Color
	bg Color
	id int16
}

func HexToColor(rrggbb string) Color {
	r, _ := strconv.ParseInt(rrggbb[1:3], 16, 0)
	g, _ := strconv.ParseInt(rrggbb[3:5], 16, 0)
	b, _ := strconv.ParseInt(rrggbb[5:7], 16, 0)
	return Color((1 << 24) + (r << 16) + (g << 8) + b)
}

func NewColorPair(fg Color, bg Color) ColorPair {
	return ColorPair{fg, bg, -1}
}

func (p ColorPair) Fg() Color {
	return p.fg
}

func (p ColorPair) Bg() Color {
	return p.bg
}

func (p ColorPair) key() int {
	return (int(p.Fg()) << 8) + int(p.Bg())
}

func (p ColorPair) is24() bool {
	return p.Fg().is24() || p.Bg().is24()
}

type ColorTheme struct {
	Fg           Color
	Bg           Color
	DarkBg       Color
	Prompt       Color
	Match        Color
	Current      Color
	CurrentMatch Color
	Spinner      Color
	Info         Color
	Cursor       Color
	Selected     Color
	Header       Color
	Border       Color
}

type BorderStyle int

const (
	BorderNone BorderStyle = iota
	BorderAround
	BorderHorizontal
)

type MouseEvent struct {
	Y      int
	X      int
	S      int
	Down   bool
	Double bool
	Mod    bool
}

var (
	ColDefault      ColorPair
	ColNormal       ColorPair
	ColPrompt       ColorPair
	ColMatch        ColorPair
	ColCurrent      ColorPair
	ColCurrentMatch ColorPair
	ColSpinner      ColorPair
	ColInfo         ColorPair
	ColCursor       ColorPair
	ColSelected     ColorPair
	ColHeader       ColorPair
	ColBorder       ColorPair
	ColUser         ColorPair
)
