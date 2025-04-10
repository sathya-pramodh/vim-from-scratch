package command

import (
	"errors"
)

type CommandType uint64

const (
	QuitCommand CommandType = iota
	UnknownCommand
)

func GetCommandFromString(command string) (CommandType, error) {
	if command == "q" {
		return QuitCommand, nil
	}
	return UnknownCommand, errors.New("unrecognized command")
}
