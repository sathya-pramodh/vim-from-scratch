package command

import (
	"errors"
)

type CommandType uint64

const (
	QuitCommand CommandType = iota
	WriteCommand
	UnknownCommand
)

func GetCommandFromString(command string) (CommandType, error) {
	switch command {
	case "q":
		return QuitCommand, nil
	case "w":
		return WriteCommand, nil
	default:
		return UnknownCommand, errors.New("unrecognized command")
	}
}
