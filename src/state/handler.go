package state

import (
	"github.com/sathya-pramodh/vim-from-scratch/src/action"
	"github.com/sathya-pramodh/vim-from-scratch/src/command"
	"github.com/sathya-pramodh/vim-from-scratch/src/mode"
)

func (t *TuiState) HandleAction(a action.Action, opts action.ActionTrigger, ch rune) {
	switch a {
	case action.NormalModeChange:
		t.Mode = mode.NormalMode
		t.CursorY, t.CursorX = t.GetPrevCursorPos(false)
		t.CommandView.Clear()
	case action.CommandModeChange:
		t.Mode = mode.CommandMode
		t.CommandView.Clear()
	case action.InsertModeChange:
		t.Mode = mode.InsertMode
		switch opts {
		case action.AppendTrigger:
			t.CursorY, t.CursorX = t.GetNextCursorPos(false)
		case action.LineEndAppendTrigger:
			y, x, err := t.GetLineEndCursorPos()
			if err != nil {
				t.WriteError(err)
				return
			}
			t.CursorX, t.CursorY = x, y
		case action.LineStartInsertTrigger:
			y, x, err := t.GetLineStartCursorPos()
			if err != nil {
				t.WriteError(err)
				return
			}
			t.CursorX, t.CursorY = x, y
		case action.NextLineInsertTrigger:
			y, x, err := t.GetLineEndCursorPos()
			if err != nil {
				t.WriteError(err)
				return
			}
			err = t.Buf.WriteToBuf('\n', x, y)
			if err != nil {
				t.WriteError(err)
				return
			}
			t.CursorX, t.CursorY = 0, y+1
		case action.PrevLineInsertTrigger:
			y, x, err := t.GetLineStartCursorPos()
			if err != nil {
				t.WriteError(err)
				return
			}
			err = t.Buf.WriteToBuf('\n', x, y)
			if err != nil {
				t.WriteError(err)
				return
			}
			t.CursorX, t.CursorY = 0, y
		}
		t.CommandView.SetStatus("-- INSERT --")
	case action.EraseLastFromCommand:
		erased := t.CommandView.EraseLastFromCommand()
		if !erased {
			t.Mode = mode.NormalMode
		}
	case action.ExecuteCommand:
		cmd, err := command.GetCommandFromString(t.CommandView.GetCommandString())
		if err != nil {
			t.WriteError(err)
			return
		}
		t.ExecuteCommand(cmd)
	case action.GotoNextLine:
		err := t.Buf.WriteToBuf('\n', t.CursorX, t.CursorY)
		if err != nil {
			t.WriteError(err)
			return
		}
		t.CursorX = 0
		t.CursorY += 1
	case action.MoveCursorLeft:
		t.MoveCursorLeft()
	case action.MoveCursorDown:
		t.MoveCursorDown()
	case action.MoveCursorUp:
		t.MoveCursorUp()
	case action.MoveCursorRight:
		t.MoveCursorRight()
	case action.InsertBackspaceChar:
		y, x := t.CursorY, t.CursorX
		t.CursorY, t.CursorX = t.GetPrevCursorPos(true)
		err := t.Buf.DeleteFromBuf(x, y)
		if err != nil {
			t.WriteError(err)
			return
		}
	case action.UnknownAction:
		switch t.Mode {
		case mode.CommandMode:
			t.CommandView.AppendToCommand(ch)
		case mode.InsertMode:
			var stringToWrite string
			if opts == action.MustTab {
				stringToWrite = "    "
			} else {
				stringToWrite = string(ch)
			}
			for _, c := range stringToWrite {
				err := t.Buf.WriteToBuf(c, t.CursorX, t.CursorY)
				if err != nil {
					t.WriteError(err)
					return
				}
			}
			t.CursorY, t.CursorX = t.GetNextCursorPos(opts == action.MustTab)
		default:
			return
		}
	}
	t.MustRefresh = true
}
