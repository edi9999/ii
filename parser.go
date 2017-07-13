package main

import (
	"github.com/gdamore/tcell"
	"github.com/edi9999/ii/core"
	"github.com/edi9999/ii/tui"
	"sync"
)

func GetChar(ev interface{}) core.Event {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		return core.Event{tui.Resize, 0, nil}

	// process mouse events:
	case *tcell.EventMouse:
		x, y := ev.Position()
		button := ev.Buttons()
		mod := ev.Modifiers() != 0
		if button&tcell.WheelDown != 0 {
			return core.Event{tui.Mouse, 0, &tui.MouseEvent{y, x, -1, false, false, mod}}
		} else if button&tcell.WheelUp != 0 {
			return core.Event{tui.Mouse, 0, &tui.MouseEvent{y, x, +1, false, false, mod}}
		}
		// process keyboard:
	case *tcell.EventKey:
		alt := (ev.Modifiers() & tcell.ModAlt) > 0
		keyfn := func(r rune) int {
			if alt {
				return tui.CtrlAltA - 'a' + int(r)
			}
			return tui.CtrlA - 'a' + int(r)
		}
		switch ev.Key() {
		case tcell.KeyCtrlA:
			return core.Event{keyfn('a'), 0, nil}
		case tcell.KeyCtrlB:
			return core.Event{keyfn('b'), 0, nil}
		case tcell.KeyCtrlC:
			return core.Event{keyfn('c'), 0, nil}
		case tcell.KeyCtrlD:
			return core.Event{keyfn('d'), 0, nil}
		case tcell.KeyCtrlE:
			return core.Event{keyfn('e'), 0, nil}
		case tcell.KeyCtrlF:
			return core.Event{keyfn('f'), 0, nil}
		case tcell.KeyCtrlG:
			return core.Event{keyfn('g'), 0, nil}
		case tcell.KeyCtrlH:
			return core.Event{keyfn('h'), 0, nil}
		case tcell.KeyCtrlI:
			return core.Event{keyfn('i'), 0, nil}
		case tcell.KeyCtrlJ:
			return core.Event{keyfn('j'), 0, nil}
		case tcell.KeyCtrlK:
			return core.Event{keyfn('k'), 0, nil}
		case tcell.KeyCtrlL:
			return core.Event{keyfn('l'), 0, nil}
		case tcell.KeyCtrlM:
			return core.Event{keyfn('m'), 0, nil}
		case tcell.KeyCtrlN:
			return core.Event{keyfn('n'), 0, nil}
		case tcell.KeyCtrlO:
			return core.Event{keyfn('o'), 0, nil}
		case tcell.KeyCtrlP:
			return core.Event{keyfn('p'), 0, nil}
		case tcell.KeyCtrlQ:
			return core.Event{keyfn('q'), 0, nil}
		case tcell.KeyCtrlR:
			return core.Event{keyfn('r'), 0, nil}
		case tcell.KeyCtrlS:
			return core.Event{keyfn('s'), 0, nil}
		case tcell.KeyCtrlT:
			return core.Event{keyfn('t'), 0, nil}
		case tcell.KeyCtrlU:
			return core.Event{keyfn('u'), 0, nil}
		case tcell.KeyCtrlV:
			return core.Event{keyfn('v'), 0, nil}
		case tcell.KeyCtrlW:
			return core.Event{keyfn('w'), 0, nil}
		case tcell.KeyCtrlX:
			return core.Event{keyfn('x'), 0, nil}
		case tcell.KeyCtrlY:
			return core.Event{keyfn('y'), 0, nil}
		case tcell.KeyCtrlZ:
			return core.Event{keyfn('z'), 0, nil}
		case tcell.KeyCtrlSpace:
			return core.Event{tui.CtrlSpace, 0, nil}
		case tcell.KeyBackspace2:
			if alt {
				return core.Event{tui.AltBS, 0, nil}
			}
			return core.Event{tui.BSpace, 0, nil}

		case tcell.KeyUp:
			return core.Event{tui.Up, 0, nil}
		case tcell.KeyDown:
			return core.Event{tui.Down, 0, nil}
		case tcell.KeyLeft:
			return core.Event{tui.Left, 0, nil}
		case tcell.KeyRight:
			return core.Event{tui.Right, 0, nil}

		case tcell.KeyHome:
			return core.Event{tui.Home, 0, nil}
		case tcell.KeyDelete:
			return core.Event{tui.Del, 0, nil}
		case tcell.KeyEnd:
			return core.Event{tui.End, 0, nil}
		case tcell.KeyPgUp:
			return core.Event{tui.PgUp, 0, nil}
		case tcell.KeyPgDn:
			return core.Event{tui.PgDn, 0, nil}

		case tcell.KeyBacktab:
			return core.Event{tui.BTab, 0, nil}

		case tcell.KeyF1:
			return core.Event{tui.F1, 0, nil}
		case tcell.KeyF2:
			return core.Event{tui.F2, 0, nil}
		case tcell.KeyF3:
			return core.Event{tui.F3, 0, nil}
		case tcell.KeyF4:
			return core.Event{tui.F4, 0, nil}
		case tcell.KeyF5:
			return core.Event{tui.F5, 0, nil}
		case tcell.KeyF6:
			return core.Event{tui.F6, 0, nil}
		case tcell.KeyF7:
			return core.Event{tui.F7, 0, nil}
		case tcell.KeyF8:
			return core.Event{tui.F8, 0, nil}
		case tcell.KeyF9:
			return core.Event{tui.F9, 0, nil}
		case tcell.KeyF10:
			return core.Event{tui.F10, 0, nil}
		case tcell.KeyF11:
			return core.Event{tui.F11, 0, nil}
		case tcell.KeyF12:
			return core.Event{tui.F12, 0, nil}

		// ev.Ch doesn't work for some reason for space:
		case tcell.KeyRune:
			r := ev.Rune()
			if alt {
				switch r {
				case ' ':
					return core.Event{tui.AltSpace, 0, nil}
				case '/':
					return core.Event{tui.AltSlash, 0, nil}
				}
				if r >= 'a' && r <= 'z' {
					return core.Event{tui.AltA + int(r) - 'a', 0, nil}
				}
				if r >= '0' && r <= '9' {
					return core.Event{tui.Alt0 + int(r) - '0', 0, nil}
				}
			}
			return core.Event{tui.Rune, r, nil}

		case tcell.KeyEsc:
			return core.Event{tui.ESC, 0, nil}

		}
	}

	return core.Event{tui.Invalid, 0, nil}
}

func ParseCommand(s tcell.Screen, commands chan core.Executer, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		ev := s.PollEvent()
		event := GetChar(ev)
		if event.Type == tui.CtrlC {
			close(commands)
			return
		}
		commands <- event
	}
}
