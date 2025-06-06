package main

import (
	"fmt"
	"os"

	"github.com/sathya-pramodh/vim-from-scratch/src/tui"
)

func main() {
	args := os.Args
	var filePath string
	if len(args) > 1 {
		filePath = args[1]
	} else {
		filePath = ""
	}
	t, err := tui.NewTui(filePath)
	if err != nil {
		fmt.Printf("%s", err)
		t.Quit()
		return
	}
	defer t.Quit()

	t.Run()
}
