package tui

import (
	"fmt"
	"os"

	"github.com/sathya-pramodh/vim-from-scratch/src/command"
	"seehuhn.de/go/ncurses"
)

type Mode uint8

const (
	NormalMode Mode = iota
	InsertMode
	VisualMode
	CommandMode
)

type Tui struct {
	win             *ncurses.Window
	buf             *Buffer
	attachedBuffers []*Buffer
	mode            Mode
	commandView     *command.CommandView
	cursorX         int
	cursorY         int
}

type Buffer struct {
	contents string
}

func NewTui() Tui {
	newBuffer := Buffer{
		contents: "VIM FROM SCRATCH",
	}
	var attachedBuffers []*Buffer
	attachedBuffers = append(attachedBuffers, &newBuffer)
	win := ncurses.Init()
	win.Erase()
	win.AddStr(newBuffer.contents)
	win.Refresh()
	commandView := command.NewCommandView(win)
	cursorY, cursorX := win.GetYX()
	return Tui{
		win:             win,
		buf:             &newBuffer,
		attachedBuffers: attachedBuffers,
		mode:            NormalMode,
		commandView:     &commandView,
		cursorX:         cursorX,
		cursorY:         cursorY,
	}
}

func (t *Tui) executeCommand(commandStr command.CommandType) {
	switch commandStr {
	case command.QuitCommand:
		t.Quit()
	}
}

func (t *Tui) Run() error {
	for {
		ch := t.win.GetCh()
		action := GetKeyAction(t.mode, ch)
		switch action {
		case ModeChange:
			prevMode := t.mode
			t.mode = GetTargetMode(t.mode, ch)
			if prevMode == CommandMode && t.mode == NormalMode {
				t.commandView.ClearCommand()
			} else if prevMode == NormalMode && t.mode == CommandMode {
				t.commandView.ClearCommand()
			}
		case AppendToCommand:
			t.commandView.AppendToCommand(ch)
		case EraseLastFromCommand:
			erased := t.commandView.EraseLastFromCommand()
			if !erased {
				t.mode = NormalMode
			}
		case ExecuteCommand:
			command, err := command.GetCommandFromString(t.commandView.GetCommandString())
			if err != nil {
				t.commandView.SetCommand(fmt.Sprintf("Error: %s", err))
				t.mode = NormalMode
				t.refresh()
				continue
			}
			t.executeCommand(command)
		case UnknownAction:
			if t.mode == CommandMode {
				t.commandView.AppendToCommand(ch)
			}
		}
		t.refresh()
	}
}

func (t *Tui) Quit() {
	t.commandView.Quit()
	os.Exit(0)
}

func (t *Tui) refresh() {
	if t.mode == NormalMode {
		t.win.Move(t.cursorY, t.cursorX)
	}
	t.commandView.Refresh(t.mode == CommandMode)
}
