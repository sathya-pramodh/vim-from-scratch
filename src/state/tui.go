package state

import (
	"fmt"
	"os"

	"github.com/sathya-pramodh/vim-from-scratch/src/buffer"
	"github.com/sathya-pramodh/vim-from-scratch/src/command"
	"github.com/sathya-pramodh/vim-from-scratch/src/mode"
	"seehuhn.de/go/ncurses"
)

type TuiState struct {
	Buf             *buffer.Buffer
	AttachedBuffers *[]*buffer.Buffer
	CursorX         int
	CursorY         int
	MaxX            int
	MaxY            int
	Mode            mode.Mode
	CommandView     *command.CommandView
	MustRefresh     bool
}

func NewTuiState(win *ncurses.Window) TuiState {
	newBuffer := buffer.Buffer{
		Contents: "",
	}
	var attachedBuffers []*buffer.Buffer
	attachedBuffers = append(attachedBuffers, &newBuffer)
	win.Erase()
	win.AddStr(newBuffer.Contents)
	win.Move(0, 0)
	win.Refresh()
	CursorY, CursorX := win.GetYX()
	MaxY, MaxX := win.GetMaxYX()
	CommandView := command.NewCommandView(win)
	return TuiState{
		Buf:             &newBuffer,
		Mode:            mode.NormalMode,
		CommandView:     &CommandView,
		AttachedBuffers: &attachedBuffers,
		CursorX:         CursorX,
		CursorY:         CursorY,
		MaxX:            MaxX,
		MaxY:            MaxY,
		MustRefresh:     false,
	}
}

func (t *TuiState) WriteError(err error) {
	t.CommandView.SetStatus(fmt.Sprintf("Error: %s", err))
	t.Mode = mode.NormalMode
	t.MustRefresh = true
}

func (t *TuiState) ExecuteCommand(cmd command.CommandType) {
	switch cmd {
	case command.QuitCommand:
		t.Quit()
		command.ExecuteQuitCommand()
	}
}

func (t *TuiState) getLineEnd(y int) (int, error) {
	x, err := t.Buf.GetLineEndX(y)
	if err != nil {
		return -1, fmt.Errorf("get line end pos: %s", err)
	}
	return x, nil
}

func (t *TuiState) getLineStart(y int) (int, error) {
	x, err := t.Buf.GetLineStartX(y)
	if err != nil {
		return -1, fmt.Errorf("get line start pos: %s", err)
	}
	return x, nil

}

func (t *TuiState) GetLineEndCursorPos() (int, int, error) {
	x, err := t.getLineEnd(t.CursorY)
	if err != nil {
		return -1, -1, err
	}
	return t.CursorY, x, nil
}

func (t *TuiState) GetLineStartCursorPos() (int, int, error) {
	x, err := t.getLineStart(t.CursorY)
	if err != nil {
		return -1, -1, err
	}
	return t.CursorY, x, nil
}

func (t *TuiState) GetPrevCursorPos(warpToPrevLine bool) (int, int) {
	if t.CursorX-1 < 0 {
		if warpToPrevLine {
			t.CursorY -= 1
			_, x, err := t.GetLineEndCursorPos()
			if err != nil {
				t.CursorY += 1
			} else {
				t.CursorX = x
			}
		}
		return t.CursorY, t.CursorX
	}
	return t.CursorY, t.CursorX - 1
}

func (t *TuiState) GetNextCursorPos(mustTab bool) (int, int) {
	var inc int
	if mustTab {
		inc = 4
	} else {
		inc = 1
	}
	if t.CursorX+inc > t.MaxX {
		return t.CursorY + 1, 1
	}
	return t.CursorY, t.CursorX + inc
}

func (t *TuiState) Quit() {
	t.CommandView.Quit()
	os.Exit(0)
}
