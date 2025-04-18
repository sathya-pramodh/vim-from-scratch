package action

import (
	"github.com/sathya-pramodh/vim-from-scratch/src/mode"
	"seehuhn.de/go/ncurses"
)

type Action uint64

const (
	InsertModeChange Action = iota
	NormalModeChange
	VisualModeChange
	CommandModeChange

	MoveCursorLeft
	MoveCursorDown
	MoveCursorUp
	MoveCursorRight
	MoveCursorNextWord
	MoveCursorNextWordEnd

	EraseLastFromCommand
	ExecuteCommand

	InsertBackspaceChar

	GotoNextLine

	UnknownAction
)

type ActionTrigger uint64

const (
	AppendTrigger ActionTrigger = iota
	LineEndAppendTrigger
	LineStartInsertTrigger
	NextLineInsertTrigger
	PrevLineInsertTrigger
	MustTab
	NoTrigger
)

func GetKeyAction(currentMode mode.Mode, ch rune) (Action, ActionTrigger) {
	switch ch {
	case ':':
		if currentMode == mode.NormalMode {
			return CommandModeChange, NoTrigger
		}
	case 'i':
		if currentMode == mode.NormalMode {
			return InsertModeChange, NoTrigger
		}
	case 'I':
		if currentMode == mode.NormalMode {
			return InsertModeChange, LineStartInsertTrigger
		}
	case 'v':
		if currentMode == mode.NormalMode {
			return VisualModeChange, NoTrigger
		}
	case 'a':
		if currentMode == mode.NormalMode {
			return InsertModeChange, AppendTrigger
		}
	case 'A':
		if currentMode == mode.NormalMode {
			return InsertModeChange, LineEndAppendTrigger
		}
	case 'o':
		if currentMode == mode.NormalMode {
			return InsertModeChange, NextLineInsertTrigger
		}
	case 'O':
		if currentMode == mode.NormalMode {
			return InsertModeChange, PrevLineInsertTrigger
		}
	case 'h':
		if currentMode == mode.NormalMode {
			return MoveCursorLeft, NoTrigger
		}
	case 'j':
		if currentMode == mode.NormalMode {
			return MoveCursorDown, NoTrigger
		}
	case 'k':
		if currentMode == mode.NormalMode {
			return MoveCursorUp, NoTrigger
		}
	case 'l':
		if currentMode == mode.NormalMode {
			return MoveCursorRight, NoTrigger
		}
	case 'w':
		if currentMode == mode.NormalMode {
			return MoveCursorNextWord, NoTrigger
		}
	case 'e':
		if currentMode == mode.NormalMode {
			return MoveCursorNextWordEnd, NoTrigger
		}
	case '\n': // <CR>
		if currentMode == mode.CommandMode {
			return ExecuteCommand, NoTrigger
		}
		if currentMode == mode.InsertMode {
			return GotoNextLine, NoTrigger
		}
	case '\t': // <Tab>
		if currentMode == mode.InsertMode {
			return UnknownAction, MustTab
		}
	case ncurses.KeyBackspace, 127: // Backspace
		if currentMode == mode.CommandMode {
			return EraseLastFromCommand, NoTrigger
		}
		if currentMode == mode.InsertMode {
			return InsertBackspaceChar, NoTrigger
		}
	case 27: // Escape
		if currentMode == mode.CommandMode || currentMode == mode.InsertMode {
			return NormalModeChange, NoTrigger
		}
	}
	return UnknownAction, NoTrigger
}
