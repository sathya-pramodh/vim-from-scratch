package tui

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
	InsertTabChar

	GotoNextLine

	UnknownAction
)

type ActionOpts struct {
	appendTrigger          bool
	lineEndAppendTrigger   bool
	lineStartInsertTrigger bool
	nextLineInsertTrigger  bool
	prevLineInsertTrigger  bool
	mustTab                bool
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
	case 'h':
		if currentMode == NormalMode {
			return MoveCursorLeft, nil
		}
	case 'j':
		if currentMode == NormalMode {
			return MoveCursorDown, nil
		}
	case 'k':
		if currentMode == NormalMode {
			return MoveCursorUp, nil
		}
	case 'l':
		if currentMode == NormalMode {
			return MoveCursorRight, nil
		}
	case '\n': // <CR>
		if currentMode == CommandMode {
			return ExecuteCommand, nil
		}
		if currentMode == InsertMode {
			return GotoNextLine, nil
		}
	case '\t': // <Tab>
		if currentMode == InsertMode {
			return UnknownAction, &ActionOpts{mustTab: true}
		}
	case 127: // Backspace
		if currentMode == CommandMode {
			return EraseLastFromCommand, nil
		}
		if currentMode == InsertMode {
			return InsertBackspaceChar, nil
		}
	case 27: // Escape
		if currentMode == CommandMode || currentMode == InsertMode {
			return NormalModeChange, nil
		}
	}
	return UnknownAction, &ActionOpts{mustTab: false}
}
