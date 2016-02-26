package syncer

import (
	"strings"

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
				l, r, err := prepair(e.Name, p)
				if err != nil {
					continue
				}
			handle:
				for _, h := range p.syncer.globalHandlers[e.Op] {
					err := h.HandleOp(l, r)
					if err == ErrReject {
						break handle
					}
				}
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

func prepair(name string, p *Pair) (l uri.Uri, r uri.Uri, err error) {
	name = strings.Replace(name, "\\", "/", -1)
	lName := p.Left.Scheme() + "://" + name
	l, err = uri.Parse(lName)
	if err != nil {
		return nil, nil, err
	}

	lTmp := p.Left.Uri()
	rTmp := p.Right.Uri()
	lTmplen := len(lTmp)
	rTmplen := len(rTmp)
	if lTmp[lTmplen-1] == '/' {
		lTmp = lTmp[0 : lTmplen-1]
	}
	if rTmp[rTmplen-1] == '/' {
		rTmp = rTmp[0 : rTmplen-1]
	}
	Uris := strings.Replace(l.Uri(), lTmp, rTmp, -1)
	r, err = uri.Parse(Uris)
	return
}
