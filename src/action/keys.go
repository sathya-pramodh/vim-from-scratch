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

	EraseLastFromCommand
	ExecuteCommand

	InsertBackspaceChar

	GotoNextLine

	UnknownAction
)

type ActionOpts struct {
	AppendTrigger          bool
	LineEndAppendTrigger   bool
	LineStartInsertTrigger bool
	NextLineInsertTrigger  bool
	PrevLineInsertTrigger  bool
	MustTab                bool
}

func GetKeyAction(currentMode mode.Mode, ch rune) (Action, *ActionOpts) {
	switch ch {
	case ':':
		if currentMode == mode.NormalMode {
			return CommandModeChange, nil
		}
	case 'i':
		if currentMode == mode.NormalMode {
			return InsertModeChange, nil
		}
	case 'I':
		if currentMode == mode.NormalMode {
			return InsertModeChange, &ActionOpts{LineStartInsertTrigger: true}
		}
	case 'v':
		if currentMode == mode.NormalMode {
			return VisualModeChange, nil
		}
	case 'a':
		if currentMode == mode.NormalMode {
			return InsertModeChange, &ActionOpts{AppendTrigger: true}
		}
	case 'A':
		if currentMode == mode.NormalMode {
			return InsertModeChange, &ActionOpts{LineEndAppendTrigger: true}
		}
	case 'o':
		if currentMode == mode.NormalMode {
			return InsertModeChange, &ActionOpts{NextLineInsertTrigger: true}
		}
	case 'O':
		if currentMode == mode.NormalMode {
			return InsertModeChange, &ActionOpts{PrevLineInsertTrigger: true}
		}
	case 'h':
		if currentMode == mode.NormalMode {
			return MoveCursorLeft, nil
		}
	case 'j':
		if currentMode == mode.NormalMode {
			return MoveCursorDown, nil
		}
	case 'k':
		if currentMode == mode.NormalMode {
			return MoveCursorUp, nil
		}
	case 'l':
		if currentMode == mode.NormalMode {
			return MoveCursorRight, nil
		}
	case '\n': // <CR>
		if currentMode == mode.CommandMode {
			return ExecuteCommand, nil
		}
		if currentMode == mode.InsertMode {
			return GotoNextLine, nil
		}
	case '\t': // <Tab>
		if currentMode == mode.InsertMode {
			return UnknownAction, &ActionOpts{MustTab: true}
		}
	case ncurses.KeyBackspace, 127: // Backspace
		if currentMode == mode.CommandMode {
			return EraseLastFromCommand, nil
		}
		if currentMode == mode.InsertMode {
			return InsertBackspaceChar, nil
		}
	case 27: // Escape
		if currentMode == mode.CommandMode || currentMode == mode.InsertMode {
			return NormalModeChange, nil
		}
	}
	return UnknownAction, &ActionOpts{MustTab: false}
}
