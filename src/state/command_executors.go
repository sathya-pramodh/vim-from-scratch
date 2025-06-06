package state

import (
	"fmt"
	"os"
)

func (t *TuiState) ExecuteQuitCommand() {
	t.Quit()
	os.Exit(0)
}

func (t *TuiState) ExecuteWriteCommand() error {
	err := t.Buf.Write()
	if err != nil {
		return fmt.Errorf("ExecuteWriteCommand: %s", err)
	}
	return nil
}
