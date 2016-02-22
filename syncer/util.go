package syncer

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func md5Hash(source ...interface{}) string {
	ctx := md5.New()
	for _, s := range source {
		switch ss := s.(type) {
		case io.Reader:
			io.Copy(ctx, ss)
		case string:
			ctx.Write([]byte(ss))
		case []byte:
			ctx.Write(ss)

		}
	}

	return hex.EncodeToString(ctx.Sum(nil))
}

func execHandler(typ ErrType, err error, h interface{}) {
	switch typ {
	case TypeInfo:
		if hi, ok := h.(InfoHandler); ok {
			hi.Info(err)
		}
	case TypeError:
		if hi, ok := h.(ErrorHandler); ok {
			hi.Error(err)
		}
	case TypeDebug:
		if hi, ok := h.(DebugHandler); ok {
			hi.Debug(err)
		}
	case TypePanic:
		if hi, ok := h.(PanicHandler); ok {
			hi.Panic(err)
		}
	case TypeFatal:
		if hi, ok := h.(FatalHandler); ok {
			hi.Fatal(err)
		}
	case TypeWarn:
		if hi, ok := h.(WarnHandler); ok {
			hi.Warn(err)
		}
	case TypeUnknown:
		if hi, ok := h.(UnknownHandler); ok {
			hi.Unknown(err)
		}
	}
}
