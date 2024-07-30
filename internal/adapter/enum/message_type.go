package enum

import (
	"database/sql/driver"
	"errors"
)

type MessageType int

const (
	User MessageType = iota + 1
	Ai
)

var MessageTypes = struct {
	User MessageType
	Ai   MessageType
}{
	User: User,
	Ai:   Ai,
}

func (mt MessageType) String() string {
	switch mt {
	case MessageTypes.User:
		return "USER"
	case MessageTypes.Ai:
		return "AI"
	default:
		return "UNKNOWN"
	}
}

func (mt *MessageType) Scan(value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("invalid type for MessageType")
	}
	*mt = MessageType(v)
	return nil
}

func (mt MessageType) Value() (driver.Value, error) {
	return int64(mt), nil
}
