package syncer

import (
	"fmt"
	"sync"

	"github.com/Felamande/filesync.v2/syncer/uri"
	fsnotify "gopkg.in/fsnotify.v1"
)

type Syncer struct {
	Pairs          []*Pair
	globalHandlers map[fsnotify.Op][]OpHandler
	errHandlers    []interface{}
}

func (s *Syncer) HandleOp(op fsnotify.Op, hs ...interface{}) {
	for _, hi := range hs {
		switch h := hi.(type) {
		case func(uri.Uri, uri.Uri) error:
			s.globalHandlers[op] = append(s.globalHandlers[op], OpHandlerFunc(h))
		case OpHandler:
			s.globalHandlers[op] = append(s.globalHandlers[op], h)
		}
	}

}

func (s *Syncer) HandleError(hs ...interface{}) {

	s.errHandlers = append(s.errHandlers, hs...)
}

func (s *Syncer) AddPair(left, right string, config *PairConfig, skip ...string) error {

	lUri, err := uri.Parse(left)
	if err != nil {
		return err
	}
	rUri, err := uri.Parse(right)
	if err != nil {
		return err
	}
	p := &Pair{
		Left:     lUri,
		Right:    rUri,
		Skip:     skip,
		handlers: make(map[fsnotify.Op][]OpHandler),
		syncer:   s,
		Config:   config,
	}
	s.Pairs = append(s.Pairs, p)
	return nil
}

func (s *Syncer) beginWatch() {
	for _, p := range s.Pairs {
		err := p.BeginWatch()
		fmt.Println("begin watch", err)

	}
}
func New() *Syncer {
	return &Syncer{
		globalHandlers: make(map[fsnotify.Op][]OpHandler, 5),
		errHandlers:    make([]interface{}, 4),
	}
}

func (s *Syncer) Run() {
	s.beginWatch()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
