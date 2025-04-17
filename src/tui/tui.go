package tui

import (
	"github.com/sathya-pramodh/vim-from-scratch/src/action"
	"github.com/sathya-pramodh/vim-from-scratch/src/mode"
	"github.com/sathya-pramodh/vim-from-scratch/src/state"
	"seehuhn.de/go/ncurses"
)

type Tui struct {
	win   *ncurses.Window
	state *state.TuiState
}

func NewTui() Tui {
	win := ncurses.Init()
	state := state.NewTuiState(win)
	return Tui{
		win:   win,
		state: &state,
	}
}

func (t *Tui) Run() {
	for {
		ch := t.win.GetCh()
		action, opts := action.GetKeyAction(t.state.Mode, ch)
		t.state.HandleAction(action, opts, ch)
		t.refresh()
	}
}

func (t *Tui) refresh() {
	if !t.state.MustRefresh {
		return
	}
	t.win.Erase()
	t.win.AddStr(t.state.Buf.Contents)
	t.win.Refresh()
	if t.state.Mode == mode.NormalMode || t.state.Mode == mode.InsertMode {
		t.win.Move(t.state.CursorY, t.state.CursorX)
	}
	t.state.CommandView.Refresh(t.state.Mode == mode.CommandMode)
	t.state.MustRefresh = false
}

func (t *Tui) Quit() {
	t.state.Quit()
}
