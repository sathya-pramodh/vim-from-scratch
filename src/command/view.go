package command

import (
	"seehuhn.de/go/ncurses"
)

type CommandView struct {
	win     *ncurses.Window
	command string
	status  string
}

func NewCommandView(win *ncurses.Window) CommandView {
	height, width := win.GetMaxYX()
	y, x := win.GetBegYX()
	subWin := ncurses.NewWin(1, width, y+height-1, x)
	return CommandView{
		win:     subWin,
		command: "",
		status:  "",
	}
}

func (c *CommandView) SetStatus(status string) {
	c.status = status
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

func (c *CommandView) Clear() {
	c.command = ""
	c.status = ""
}

func (c *CommandView) GetCommandString() string {
	return c.command
}

func (c *CommandView) Refresh(inCommandMode bool) {
	c.win.Erase()
	if inCommandMode {
		c.win.AddStr(":" + c.command)
	} else if c.status != "" {
		c.win.AddStr(c.status)
	}
	c.win.Refresh()
}

func (c *CommandView) Quit() {
	ncurses.EndWin()
}
