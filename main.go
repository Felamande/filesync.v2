package main

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Felamande/filesync.v2/syncer"
	"github.com/Felamande/filesync.v2/syncer/uri"
	"github.com/qiniu/log"
	fsnotify "gopkg.in/fsnotify.v1"
)

func main() {
	s := syncer.New()
	s.AddPair("local://C:/users/kirigiri/pictures/", "local://D:/Dev/test2", &syncer.PairConfig{})
	// s.AddPair("local://D:/music", "local://D:/Dev/test2", &syncer.PairConfig{})
	s.HandleOp(fsnotify.Create, HandleAddNewWatch, HandleCreate)
	s.HandleOp(fsnotify.Write, HandleWrite)
	s.HandleError(new(QiniuLogger).Init(os.Stdout))
	s.Run()
}

func HandleWrite(ctx syncer.Context, l uri.Uri, r uri.Uri) error {
	if l.ModTime().Sub(r.ModTime()) < 0 {
		return ctx.Finish()
	}
	var (
		reader io.ReadCloser
		writer io.WriteCloser
		err    error
	)

	for {
		reader, err = l.OpenRead()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 20)
	}
	for {
		writer, err = r.OpenWrite()
		if err == nil {
			break
		}
		time.Sleep(time.Minute * 10)
	}
	defer reader.Close()
	defer writer.Close()
	_, err = io.Copy(writer, reader)
	if err != nil {
		ctx.EmitLog(syncer.TypeError, err)
	}
	ctx.EmitLog(syncer.TypeInfo, "write to ", r.Uri())
	return err
}

func HandleAddNewWatch(ctx syncer.Context, l uri.Uri, r uri.Uri) error {

	if filepath.Base(l.Abs()) == ".git" {
		return ctx.Finish()
	}
	if !l.IsDir() {
		return nil
	}
	err := ctx.AddWatch(l)
	if err != nil {
		ctx.EmitLog(syncer.TypeError, err)
	}

	return nil
}

func HandleCreate(ctx syncer.Context, l uri.Uri, r uri.Uri) error {
	if l.ModTime().Sub(r.ModTime()) < 0 {
		return ctx.Finish()
	}
	var err error
	for {
		err = r.Create(l.IsDir(), l.Mode())
		if err == nil {
			break
		}
		time.Sleep(time.Minute * 1)
	}
	ctx.EmitLog(syncer.TypeInfo, "create ", r.Uri())
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
	default:
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
