package command

import (
	"seehuhn.de/go/ncurses"
)

type CommandView struct {
	win     *ncurses.Window
	command string
}

func NewCommandView(win *ncurses.Window) CommandView {
	height, width := win.GetMaxYX()
	y, x := win.GetBegYX()
	subWin := ncurses.NewWin(1, width, y+height-1, x)
	return CommandView{
		win:     subWin,
		command: "",
	}
}

func (c *CommandView) SetCommand(command string) {
	c.command = command
}

func (c *CommandView) AppendToCommand(ch rune) {
	c.command += string(ch)
}

func (c *CommandView) EraseLastFromCommand() bool {
	if c.command == "" {
		return false
	}
	c.command = c.command[0 : len(c.command)-1]
	return true
}

func (c *CommandView) ClearCommand() {
	c.command = ""
}

func (c *CommandView) GetCommandString() string {
	return c.command
}

func (c *CommandView) Refresh(inCommandMode bool) {
	c.win.Erase()
	if inCommandMode {
		c.win.AddStr(":" + c.command)
	} else if c.command != "" {
		c.win.AddStr(c.command)
	}
	c.win.Refresh()
}

func (c *CommandView) Quit() {
	ncurses.EndWin()
}
