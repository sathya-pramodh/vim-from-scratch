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
	maxX            int
	maxY            int
}

func NewTui() Tui {
	newBuffer := Buffer{
		contents: "",
	}
	var attachedBuffers []*Buffer
	attachedBuffers = append(attachedBuffers, &newBuffer)
	win := ncurses.Init()
	win.Erase()
	win.AddStr(newBuffer.contents)
	win.Move(0, 0)
	win.Refresh()
	commandView := command.NewCommandView(win)
	cursorY, cursorX := win.GetYX()
	maxY, maxX := win.GetMaxYX()
	return Tui{
		win:             win,
		buf:             &newBuffer,
		attachedBuffers: attachedBuffers,
		mode:            NormalMode,
		commandView:     &commandView,
		cursorX:         cursorX,
		cursorY:         cursorY,
		maxX:            maxX,
		maxY:            maxY,
	}
}

func (t *Tui) Run() {
	for {
		ch := t.win.GetCh()
		action, opts := GetKeyAction(t.mode, ch)
		switch action {
		case NormalModeChange:
			t.mode = NormalMode
			t.cursorY, t.cursorX = t.getPrevCursorPos()
			t.commandView.Clear()
		case CommandModeChange:
			t.mode = CommandMode
			t.commandView.Clear()
		case InsertModeChange:
			t.mode = InsertMode
			if opts != nil {
				if opts.appendTrigger {
					t.cursorY, t.cursorX = t.getNextCursorPos()
				}
				if opts.lineEndAppendTrigger {
					y, x, err := t.getLineEndCursorPos()
					if err != nil {
						t.writeError(err)
						continue
					}
					t.cursorX, t.cursorY = x, y
				}
				if opts.lineStartInsertTrigger {
					y, x, err := t.getLineStartCursorPos()
					if err != nil {
						t.writeError(err)
						continue
					}
					t.cursorX, t.cursorY = x, y
				}
				if opts.nextLineInsertTrigger {
					y, x, err := t.getLineEndCursorPos()
					if err != nil {
						t.writeError(err)
						continue
					}
					err = t.buf.WriteToBuf('\n', x, y)
					if err != nil {
						t.writeError(err)
						continue
					}
					t.cursorX, t.cursorY = 0, y+1
				}
				if opts.prevLineInsertTrigger {
					y, x, err := t.getLineStartCursorPos()
					if err != nil {
						t.writeError(err)
						continue
					}
					err = t.buf.WriteToBuf('\n', x, y)
					if err != nil {
						t.writeError(err)
						continue
					}
					t.cursorX, t.cursorY = 0, y
				}
			}
			t.commandView.SetStatus("--INSERT--")
		case EraseLastFromCommand:
			erased := t.commandView.EraseLastFromCommand()
			if !erased {
				t.mode = NormalMode
			}
		case ExecuteCommand:
			cmd, err := command.GetCommandFromString(t.commandView.GetCommandString())
			if err != nil {
				t.writeError(err)
				continue
			}
			t.executeCommand(cmd)
		case GotoNextLine:
			err := t.buf.WriteToBuf(ch, t.cursorX, t.cursorY)
			if err != nil {
				t.writeError(err)
				continue
			}
			t.cursorX = 0
			t.cursorY += 1
		case UnknownAction:
			switch t.mode {
			case CommandMode:
				t.commandView.AppendToCommand(ch)
			case InsertMode:
				err := t.buf.WriteToBuf(ch, t.cursorX, t.cursorY)
				if err != nil {
					t.writeError(err)
					continue
				}
				t.cursorY, t.cursorX = t.getNextCursorPos()
			default:
				continue
			}
		}
		t.refresh()
	}
}

func (t *Tui) writeError(err error) {
	t.commandView.SetStatus(fmt.Sprintf("Error: %s", err))
	t.mode = NormalMode
	t.refresh()
}

func (t *Tui) executeCommand(cmd command.CommandType) {
	switch cmd {
	case command.QuitCommand:
		t.Quit()
		command.ExecuteQuitCommand()
	}
}

func (t *Tui) getNextCursorPos() (int, int) {
	if t.cursorX+1 > t.maxX {
		return t.cursorY + 1, 1
	}
	return t.cursorY, t.cursorX + 1
}

func (t *Tui) getLineEnd(y int) (int, error) {
	x, err := t.buf.GetLineEndX(y)
	if err != nil {
		return -1, fmt.Errorf("get line end pos: %s", err)
	}
	return x, nil
}

func (t *Tui) getLineStart(y int) (int, error) {
	x, err := t.buf.GetLineStartX(y)
	if err != nil {
		return -1, fmt.Errorf("get line start pos: %s", err)
	}
	return x, nil

}

func (t *Tui) getLineEndCursorPos() (int, int, error) {
	x, err := t.getLineEnd(t.cursorY)
	if err != nil {
		return -1, -1, fmt.Errorf("get line end pos: %s", err)
	}
	return t.cursorY, x, nil
}

func (t *Tui) getLineStartCursorPos() (int, int, error) {
	x, err := t.getLineStart(t.cursorY)
	if err != nil {
		return -1, -1, fmt.Errorf("get line start pos: %s", err)
	}
	return t.cursorY, x, nil
}

func (t *Tui) getPrevCursorPos() (int, int) {
	if t.cursorX-1 < 0 {
		return t.cursorY, t.cursorX
	}
	return t.cursorY, t.cursorX - 1
}

func (t *Tui) Quit() {
	t.commandView.Quit()
	os.Exit(0)
}

func (t *Tui) refresh() {
	t.win.Erase()
	t.win.AddStr(t.buf.contents)
	t.win.Refresh()
	if t.mode == NormalMode || t.mode == InsertMode {
		t.win.Move(t.cursorY, t.cursorX)
	}
	t.commandView.Refresh(t.mode == CommandMode)
}
