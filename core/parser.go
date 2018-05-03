package core

import (
	"github.com/edi9999/ii/events"
	"github.com/gdamore/tcell"
	"sync"
)

func GetChar(ev interface{}) Event {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		return Event{events.Resize, 0, nil}

	case *tcell.EventMouse:
		x, y := ev.Position()
		button := ev.Buttons()
		mod := ev.Modifiers() != 0
		if button&tcell.WheelDown != 0 {
			return Event{events.Mouse, 0, &events.MouseEvent{y, x, -1, false, false, mod}}
		} else if button&tcell.WheelUp != 0 {
			return Event{events.Mouse, 0, &events.MouseEvent{y, x, +1, false, false, mod}}
		}
	case *tcell.EventKey:
		alt := (ev.Modifiers() & tcell.ModAlt) > 0
		keyfn := func(r rune) int {
			if alt {
				return events.CtrlAltA - 'a' + int(r)
			}
			return events.CtrlA - 'a' + int(r)
		}
		switch ev.Key() {
		case tcell.KeyCtrlA:
			return Event{keyfn('a'), 0, nil}
		case tcell.KeyCtrlB:
			return Event{keyfn('b'), 0, nil}
		case tcell.KeyCtrlC:
			return Event{keyfn('c'), 0, nil}
		case tcell.KeyCtrlD:
			return Event{keyfn('d'), 0, nil}
		case tcell.KeyCtrlE:
			return Event{keyfn('e'), 0, nil}
		case tcell.KeyCtrlF:
			return Event{keyfn('f'), 0, nil}
		case tcell.KeyCtrlG:
			return Event{keyfn('g'), 0, nil}
		case tcell.KeyCtrlH:
			return Event{keyfn('h'), 0, nil}
		case tcell.KeyCtrlI:
			return Event{keyfn('i'), 0, nil}
		case tcell.KeyCtrlJ:
			return Event{keyfn('j'), 0, nil}
		case tcell.KeyCtrlK:
			return Event{keyfn('k'), 0, nil}
		case tcell.KeyCtrlL:
			return Event{keyfn('l'), 0, nil}
		case tcell.KeyCtrlM:
			return Event{keyfn('m'), 0, nil}
		case tcell.KeyCtrlN:
			return Event{keyfn('n'), 0, nil}
		case tcell.KeyCtrlO:
			return Event{keyfn('o'), 0, nil}
		case tcell.KeyCtrlP:
			return Event{keyfn('p'), 0, nil}
		case tcell.KeyCtrlQ:
			return Event{keyfn('q'), 0, nil}
		case tcell.KeyCtrlR:
			return Event{keyfn('r'), 0, nil}
		case tcell.KeyCtrlS:
			return Event{keyfn('s'), 0, nil}
		case tcell.KeyCtrlT:
			return Event{keyfn('t'), 0, nil}
		case tcell.KeyCtrlU:
			return Event{keyfn('u'), 0, nil}
		case tcell.KeyCtrlV:
			return Event{keyfn('v'), 0, nil}
		case tcell.KeyCtrlW:
			return Event{keyfn('w'), 0, nil}
		case tcell.KeyCtrlX:
			return Event{keyfn('x'), 0, nil}
		case tcell.KeyCtrlY:
			return Event{keyfn('y'), 0, nil}
		case tcell.KeyCtrlZ:
			return Event{keyfn('z'), 0, nil}
		case tcell.KeyCtrlSpace:
			return Event{events.CtrlSpace, 0, nil}
		case tcell.KeyBackspace2:
			if alt {
				return Event{events.AltBS, 0, nil}
			}
			return Event{events.BSpace, 0, nil}

		case tcell.KeyUp:
			return Event{events.Up, 0, nil}
		case tcell.KeyDown:
			return Event{events.Down, 0, nil}
		case tcell.KeyLeft:
			return Event{events.Left, 0, nil}
		case tcell.KeyRight:
			return Event{events.Right, 0, nil}

		case tcell.KeyHome:
			return Event{events.Home, 0, nil}
		case tcell.KeyDelete:
			return Event{events.Del, 0, nil}
		case tcell.KeyEnd:
			return Event{events.End, 0, nil}
		case tcell.KeyPgUp:
			return Event{events.PgUp, 0, nil}
		case tcell.KeyPgDn:
			return Event{events.PgDn, 0, nil}

		case tcell.KeyBacktab:
			return Event{events.BTab, 0, nil}

		case tcell.KeyF1:
			return Event{events.F1, 0, nil}
		case tcell.KeyF2:
			return Event{events.F2, 0, nil}
		case tcell.KeyF3:
			return Event{events.F3, 0, nil}
		case tcell.KeyF4:
			return Event{events.F4, 0, nil}
		case tcell.KeyF5:
			return Event{events.F5, 0, nil}
		case tcell.KeyF6:
			return Event{events.F6, 0, nil}
		case tcell.KeyF7:
			return Event{events.F7, 0, nil}
		case tcell.KeyF8:
			return Event{events.F8, 0, nil}
		case tcell.KeyF9:
			return Event{events.F9, 0, nil}
		case tcell.KeyF10:
			return Event{events.F10, 0, nil}
		case tcell.KeyF11:
			return Event{events.F11, 0, nil}
		case tcell.KeyF12:
			return Event{events.F12, 0, nil}

		// ev.Ch doesn't work for some reason for space:
		case tcell.KeyRune:
			r := ev.Rune()
			if alt {
				switch r {
				case ' ':
					return Event{events.AltSpace, 0, nil}
				case '/':
					return Event{events.AltSlash, 0, nil}
				}
				if r >= 'a' && r <= 'z' {
					return Event{events.AltA + int(r) - 'a', 0, nil}
				}
				if r >= '0' && r <= '9' {
					return Event{events.Alt0 + int(r) - '0', 0, nil}
				}
			}
			return Event{events.Rune, r, nil}

		case tcell.KeyEsc:
			return Event{events.ESC, 0, nil}

		}
	}

	return Event{events.Invalid, 0, nil}
}

func InputParser(s tcell.Screen, commands chan Executer, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		ev := s.PollEvent()
		event := GetChar(ev)
		if event.Type == events.CtrlC {
			close(commands)
			return
		}
		commands <- event
	}
}
