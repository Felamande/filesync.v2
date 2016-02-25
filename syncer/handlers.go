package syncer

import (
	"fmt"

	"github.com/Felamande/filesync.v2/syncer/uri"
)

func DefaultSimpleErrHandlers(err error) {
	fmt.Println(err)
}

type InfoHandler interface {
	Info(error)
}
type ErrorHandler interface {
	Error(error)
}
type DebugHandler interface {
	Debug(error)
}
type PanicHandler interface {
	Panic(error)
}
type FatalHandler interface {
	Fatal(error)
}
type WarnHandler interface {
	Warn(error)
}
type UnknownHandler interface {
	Unknown(error)
}

type OpHandler interface {
	HandleOp(l uri.Uri, r uri.Uri) error
}
