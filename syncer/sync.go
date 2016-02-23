package syncer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Felamande/filesync.v2/syncer/uri"
	fsnotify "gopkg.in/fsnotify.v1"
)

type Syncer struct {
	Pairs          []*Pair
	globalHandlers map[fsnotify.Op][]OpHandler
	errHandlers    []interface{}
}

func (s *Syncer) HandleOp(op fsnotify.Op, hs ...OpHandler) {

	s.globalHandlers[op] = append(s.globalHandlers[op], hs...)
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
