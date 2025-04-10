package tui

import (
	"seehuhn.de/go/ncurses"
)

type Action uint64

const (
	InsertModeChange Action = iota
	NormalModeChange
	VisualModeChange
	CommandModeChange

	AppendToCommand
	EraseLastFromCommand
	ExecuteCommand

	GotoNextLine

	UnknownAction
)

type ActionOpts struct {
	appendTrigger          bool
	lineEndAppendTrigger   bool
	lineStartInsertTrigger bool
	nextLineInsertTrigger  bool
	prevLineInsertTrigger  bool
}

func GetKeyAction(currentMode Mode, ch rune) (Action, *ActionOpts) {
	switch ch {
	case ':':
		if currentMode == NormalMode {
			return CommandModeChange, nil
		}
	case 'i':
		if currentMode == NormalMode {
			return InsertModeChange, nil
		}
	case 'I':
		if currentMode == NormalMode {
			return InsertModeChange, &ActionOpts{lineStartInsertTrigger: true}
		}
	case 'v':
		if currentMode == NormalMode {
			return VisualModeChange, nil
		}
	case 'a':
		if currentMode == NormalMode {
			return InsertModeChange, &ActionOpts{appendTrigger: true}
		}
	case 'A':
		if currentMode == NormalMode {
			return InsertModeChange, &ActionOpts{lineEndAppendTrigger: true}
		}
	case 'o':
		if currentMode == NormalMode {
			return InsertModeChange, &ActionOpts{nextLineInsertTrigger: true}
		}
	case 'O':
		if currentMode == NormalMode {
			return InsertModeChange, &ActionOpts{prevLineInsertTrigger: true}
		}
	case 'q':
		if currentMode == CommandMode {
			return AppendToCommand, nil
		}
	case '\n':
		if currentMode == CommandMode {
			return ExecuteCommand, nil
		}
		if currentMode == InsertMode {
			return GotoNextLine, nil
		}
	case ncurses.KeyBackspace:
		if currentMode == CommandMode {
			return EraseLastFromCommand, nil
		}
	case rune(27):
		if currentMode == CommandMode || currentMode == InsertMode {
			return NormalModeChange, nil
		}
	}
	return UnknownAction, nil
}
