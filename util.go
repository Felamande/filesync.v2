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
