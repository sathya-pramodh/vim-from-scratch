package mode

type Mode uint8

const (
	NormalMode Mode = iota
	InsertMode
	VisualMode
	CommandMode
)
