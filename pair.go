package syncer

import (
	"errors"

	"github.com/Felamande/syncer/uri"
	"gopkg.in/fsnotify.v1"
)

type PairConfig struct {
}

type Pair struct {
	Left  uri.Uri
	Right uri.Uri

	watcher  *fsnotify.Watcher
	progress chan int64
	syncer   *Syncer
	handlers map[fsnotify.Op][]OpHandler

	Config *PairConfig
}

func (p *Pair) clone() Pair {
	return Pair{
		Left:   p.Left,
		Right:  p.Right,
		Config: &PairConfig{},
	}
}

func (p *Pair) AddWatch(u uri.Uri) error {
	if p.watcher == nil {
		return errors.New("nil watcher")
	}
	return p.watcher.Add(u.Abs())
}
