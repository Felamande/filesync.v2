package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Felamande/syncer"
	"github.com/Felamande/syncer/uri"
	"github.com/qiniu/log"
	fsnotify "gopkg.in/fsnotify.v1"
)

func main() {
	s := syncer.New()
	s.AddPair("local://C:/users/kirigiri/pictures", "local://D:/Dev/test2", &syncer.PairConfig{})
	s.AddPair("local://D:/music", "local://D:/Dev/test2", &syncer.PairConfig{})
	s.HandleOp(fsnotify.Create, HandleCreate)
	s.HandleOp(fsnotify.Create, HandleCreate2)
	s.HandleError(new(QiniuLogger).Init(os.Stdout))
	// s.HandleError()
	s.Run()
}

func HandleCreate(ctx syncer.Context, l uri.Uri, r uri.Uri) error {
	if strings.Contains(l.Uri(), "pixiv") {
		ctx.EmitErr(syncer.TypeInfo, "got pixiv")
		ctx.EmitErr(syncer.TypeError, "got pixiv")
		ctx.EmitErr(syncer.TypeFatal, "got pixiv")
	}
	// fmt.Println(l.Uri(), r.Uri())
	if filepath.Base(l.Abs()) == ".git" {
		return syncer.ErrReject
	}
	return nil
}

func HandleCreate2(pair *syncer.Pair, l uri.Uri, r uri.Uri) error {
	// fmt.Println("invoke me")
	if l.IsDir() {
		pair.AddWatch(l)
	}

	return nil
}

type QiniuLogger struct {
	writer io.Writer
	logger *log.Logger
}

func (l *QiniuLogger) Init(w interface{}) *QiniuLogger {
	switch writer := w.(type) {
	case io.Writer:
		l.writer = writer
	case string:
		fd, err := os.OpenFile(writer, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		l.writer = fd
	}

	if l.writer == nil {
		l.writer = os.Stdout
	}

	l.logger = log.New(l.writer, "[filesync]", log.LstdFlags|log.Llevel)
	return l
}

func (l *QiniuLogger) Info(err error) {
	l.logger.Info(err)
}

func (l *QiniuLogger) Error(err error) {
	l.logger.Error(err)
}
