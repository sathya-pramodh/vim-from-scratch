package main

import (
	"fmt"
	"os"

	"github.com/sathya-pramodh/vim-from-scratch/src/tui"
)

func main() {
	t := tui.NewTui()
	defer t.Quit()

	err := t.Run()
	if err != nil {
		fmt.Printf("unable to run tui: %s\n", err)
		os.Exit(1)
	}
}
