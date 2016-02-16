package syncer

import (
	"fmt"
)

func DefaultErrHandler(err error) {
	fmt.Println(err)
}

func QiniuLogger(err error) {

}
