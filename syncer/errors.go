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
	TypeUnknown
)

type Error struct {
	typ ErrType
	err interface{}
}

func (e *Error) Type() ErrType {
	return e.typ
}

func (e *Error) Raw() interface{} {
	return e.err
}

func (e *Error) Error() string {
	switch err := e.err.(type) {
	case []interface{}:
		return fmt.Sprint(err...)
	case error:
		return err.Error()
	default:
		return ""
	}
}
