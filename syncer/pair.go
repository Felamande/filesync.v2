package syncer

import (
	"github.com/Felamande/filesync.v2/syncer/rnotify"
	"github.com/Felamande/filesync.v2/syncer/uri"
	"gopkg.in/fsnotify.v1"
)

type PairConfig struct {
}

type Pair struct {
	Left     uri.Uri
	Right    uri.Uri
	Skip     []string
	rwatcher *rnotify.Watcher

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
func (p *Pair) BeginWatch() (err error) {

	p.rwatcher, err = rnotify.NewWatcher(p.Left.Abs())
	if err != nil {
		return
	}
	events, errors, err := p.rwatcher.Skip(p.Skip...).Start()
	if err != nil {
		return
	}
	go func() {
		for {
			select {
			case e := <-events:
				println(e.Name)
			}
		}
	}()
	go func() {
		for {
			select {
			case e := <-errors:
				println(e.Error())
			}
		}
	}()
	return
}

func (p *Pair) AddToWatcher(u uri.Uri) error {
	return p.rwatcher.Add(u.Abs())
}
