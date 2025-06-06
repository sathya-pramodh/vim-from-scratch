package state

import "fmt"

func (t *TuiState) MoveCursorLeft() {
	if t.CursorX-1 < 0 {
		return
	}
	t.CursorX -= 1
}

func (t *TuiState) MoveCursorDown() {
	if t.CursorY+1 > t.MaxY {
		return
	}
	t.CursorY += 1
	_, xend, err := t.GetLineEndCursorPos()
	if err != nil {
		t.CursorY -= 1
		return
	}
	_, xstart, err := t.GetLineStartCursorPos()
	if err != nil {
		t.CursorY -= 1
		return
	}
	if t.CursorX > xend {
		t.CursorX = xend
	}
	if t.CursorX < xstart {
		t.CursorX = xstart
	}
}

func (t *TuiState) MoveCursorUp() {
	if t.CursorY-1 < 0 {
		return
	}
	t.CursorY -= 1
	_, xend, err := t.GetLineEndCursorPos()
	if err != nil {
		t.CursorY += 1
		return
	}
	_, xstart, err := t.GetLineStartCursorPos()
	if err != nil {
		t.CursorY += 1
		return
	}
	if t.CursorX > xend {
		t.CursorX = xend
	}
	if t.CursorX < xstart {
		t.CursorX = xstart
	}
}

func (t *TuiState) MoveCursorRight() {
	_, xend, _ := t.GetLineEndCursorPos()
	if t.CursorX+1 >= xend {
		return
	}
	t.CursorX += 1
}

func (t *TuiState) MoveCursorNextWord() error {
	x, y, err := t.Buf.GetNextWordPos(t.CursorX, t.CursorY)
	if err != nil {
		return fmt.Errorf("MoveCursorNextWord: %s", err)
	}
	t.CursorX, t.CursorY = x, y
	return nil
}

func (t *TuiState) MoveCursorNextWordEnd() error {
	x, y, err := t.Buf.GetNextWordEndPos(t.CursorX, t.CursorY)
	if err != nil {
		return fmt.Errorf("MoveCursorNextWordEnd: %s", err)
	}
	t.CursorX, t.CursorY = x, y
	return nil
}
