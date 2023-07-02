package types

type Action = byte

const (
	ActionNone Action = iota
	ActionNext
	ActionPrev
)
