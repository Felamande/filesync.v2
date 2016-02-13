package syncer

import (
	"github.com/Felamande/filesync.v2/uri"
	fsnotify "gopkg.in/fsnotify.v1"
)

type Handler func(p *Pair) error

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
	Pairs    []*Pair
	pmap     map[string]*Pair
	handlers map[interface{}][]Handler

	cache cacher
	errs  chan error
	msg   chan message
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

		}
	}
}

func (s *Syncer) HandleOp(op fsnotify.Op, h ...Handler) {
	s.handlers[op] = append(s.handlers[op], h...)
}

func (s *Syncer) HandleUri(uri string, h ...Handler) {
	// if !s.watching[uri] {
	// 	return
	// }
	s.handlers[uri] = append(s.handlers[uri], h...)
}

func (s *Syncer) NewPair(left, right string) (*Pair, error) {
	if p, ok := s.pmap[md5Hash(left, right)]; ok {
		return p, nil
	}
	lUri, err := s.cache.Get(left)
	if err != nil {
		return nil, err
	}
	rUri, err := s.cache.Get(right)
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
