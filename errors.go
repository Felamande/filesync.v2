package syncer

import (
	"errors"
	"fmt"
)

var ErrReject = errors.New("reject")

type ErrType int

const (
	TypeInfo ErrType = iota
	TypeError
	TypeDebug
	TypePanic
	TypeFatal
	TypeWarn
	TypeAll
	TypeUnknown
)

type Record []interface{}

type Error struct {
	typ ErrType
	err interface{}
}

func (e *Error) Error() string {
	switch err := e.err.(type) {
	case Record:
	case []interface{}:
		return fmt.Sprint(err...)
	case error:
		return err.Error()
	default:
		return ""
	}
	return ""
}
