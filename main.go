package main

import (
	"github.com/sathya-pramodh/vim-from-scratch/src/tui"
)

func main() {
	t := tui.NewTui()
	defer t.Quit()

	t.Run()
}
