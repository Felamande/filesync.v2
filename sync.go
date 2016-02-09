package syncer

import (
	"github.com/Felamande/filesync.v2/uri"
	fsnotify "gopkg.in/fsnotify.v1"
)

type Handler func(p *Pair) error

type Pair struct {
	Left    uri.Uri
	Right   uri.Uri
	written bool
	watcher *fsnotify.Watcher
}

type Syncer struct {
	Pairs []*Pair
	mux   map[fsnotify.Op][]Handler
	serv  map[string]Pair
}

func (p *Pair) Written() bool {
	return p.written
}
