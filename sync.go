package syncer

import (
	"github.com/Felamande/filesync.v2/uri"
	fsnotify "gopkg.in/fsnotify.v1"
)

type Pair struct {
	Left  uri.Uri
	Right uri.Uri

	watcher  *fsnotify.Watcher
	progress chan int64
	syncer   *Syncer
	hash     string
}

func (p *Pair) Hash() string {
	if p.hash != "" {
		return p.hash
	}
	p.hash = md5Hash(p.Left.Uri(), p.Right.Uri())
	return p.hash
}

func (p *Pair) Syncer() *Syncer {
	return p.syncer
}

type message struct {
	p  *Pair
	op fsnotify.Op
}

type Syncer struct {
	Pairs       []*Pair
	pmap        map[string]*Pair
	handlers    map[fsnotify.Op][]func(*Pair) error
	errHandlers []func(error)
	errs        chan error
	msg         chan message
}

func (s *Syncer) Mux() {
	for {
		select {
		case msg := <-s.msg:
			go func(m message) {
				for _, h := range s.handlers[m.op] {
					if h != nil {
						h(m.p)
					}
				}
			}(msg)
		case err := <-s.errs:
			go func(e error) {
				if len(s.errHandlers) == 0 {
					DefaultErrHandler(e)
				}
				for _, h := range s.errHandlers {
					h(e)
				}
			}(err)

		}
	}
}

func (s *Syncer) HandleOp(op fsnotify.Op, h ...func(p *Pair) error) {
	s.handlers[op] = append(s.handlers[op], h...)
}

func (s *Syncer) HandleError(h ...func(error)) {
	s.errHandlers = append(s.errHandlers, h...)
}

func (s *Syncer) NewPair(left, right string) (*Pair, error) {
	if p, ok := s.pmap[md5Hash(left, right)]; ok {
		return p, nil
	}
	lUri, err := uri.Parse(left)
	if err != nil {
		return nil, err
	}
	rUri, err := uri.Parse(right)
	if err != nil {
		return nil, err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	p := &Pair{
		Left:     lUri,
		Right:    rUri,
		watcher:  watcher,
		progress: make(chan int64),
		syncer:   s,
	}
	s.pmap[p.Hash()] = p
	return p, nil
}
