package tui

import (
	"seehuhn.de/go/ncurses"
)

type Action uint64

const (
	ModeChange Action = iota

	AppendToCommand
	EraseLastFromCommand
	ExecuteCommand

	UnknownAction
)

func GetKeyAction(currentMode Mode, ch rune) Action {
	switch ch {
	case ':':
		if currentMode == NormalMode {
			return ModeChange
		}
	case 'q':
		if currentMode == CommandMode {
			return AppendToCommand
		}
	case '\n':
		if currentMode == CommandMode {
			return ExecuteCommand
		}
	case ncurses.KeyBackspace:
		if currentMode == CommandMode {
			return EraseLastFromCommand
		}
	case rune(27):
		if currentMode == CommandMode {
			return ModeChange
		}
	}
	return UnknownAction
}

func GetTargetMode(currentMode Mode, ch rune) Mode {
	switch ch {
	case ':':
		if currentMode == NormalMode {
			return CommandMode
		}
	case 'v':
		if currentMode == NormalMode {
			return VisualMode
		}
	case 'i':
		if currentMode == NormalMode {
			return InsertMode
		}
	case rune(27):
		if currentMode == CommandMode {
			return NormalMode
		}
	}
	return currentMode
}
