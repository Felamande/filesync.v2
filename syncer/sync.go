package syncer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Felamande/filesync.v2/syncer/uri"
	fsnotify "gopkg.in/fsnotify.v1"
)

type Message struct {
	Name string
	P    *Pair
	Op   fsnotify.Op
}

type Syncer struct {
	Pairs          []*Pair
	globalHandlers map[fsnotify.Op][]OpHandler
	errHandlers    []interface{}
	errs           chan error
	msg            chan Message
	addPair        chan *Pair
}

func (s *Syncer) loopMsg() {
	for {
		select {
		case msg := <-s.msg:
			go func(m Message) {

				handlers, exist := s.globalHandlers[m.Op]
				if !exist {
					return
				}
				if len(handlers) == 0 {
					return
				}
				l, r, err := prepair(m.Name, m.P)
				if err != nil {
					go func() { s.errs <- err }()
					return
				}
				for _, h := range handlers {
					if h == nil {
						continue
					}
					err := h.HandleOp(Context{m.P, s}, l, r)
					if err == ErrReject {
						return
					}
				}
			}(msg)
		}
	}
}

func (s *Syncer) loopErr() {
	for {
		select {
		case err := <-s.errs:
			go func(e error) {
				if e == nil {
					return
				}
				if E, ok := e.(*Error); ok {
					for _, h := range s.errHandlers {
						if h == nil {
							continue
						}
						execHandler(E.typ, e, h)
					}
				} else {
					go func() { s.errs <- &Error{TypeUnknown, e} }()
				}

			}(err)
		}
	}
}

func (s *Syncer) HandleOp(op fsnotify.Op, hs ...interface{}) {
	for _, h := range hs {
		switch handler := h.(type) {
		case func(Context, uri.Uri, uri.Uri) error:
			s.globalHandlers[op] = append(s.globalHandlers[op], OpHandlerFunc(handler))

		default:
			if iHandler, ok := h.(OpHandler); ok {
				s.globalHandlers[op] = append(s.globalHandlers[op], iHandler)
			}
		}

	}
}

func (s *Syncer) HandleError(hs ...interface{}) {

	s.errHandlers = append(s.errHandlers, hs...)
}

func (s *Syncer) AddPair(left, right string, config *PairConfig) error {

	lUri, err := uri.Parse(left)
	if err != nil {
		return err
	}
	rUri, err := uri.Parse(right)
	if err != nil {
		return err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	p := &Pair{
		Left:     lUri,
		Right:    rUri,
		watcher:  watcher,
		progress: make(chan int64),
		handlers: make(map[fsnotify.Op][]OpHandler),
		syncer:   s,
		Config:   config,
	}
	s.Pairs = append(s.Pairs, p)
	go func(pair *Pair) {
		s.addPair <- pair
		fmt.Println("add pair", pair.Left.Uri())
	}(p)
	return nil
}

func (s *Syncer) emitErr(typ ErrType, v ...interface{}) {
	go func() { s.errs <- &Error{typ, v} }()
}

func (s *Syncer) BeginWatch() {
	for {
		select {
		case p := <-s.addPair:
			p.Left.Walk(&localVisitor{p, s})
			go func(pair *Pair) {
				for {
					select {
					case event := <-pair.watcher.Events:
						go func(e fsnotify.Event) {
							s.msg <- Message{e.Name, pair, e.Op}
						}(event)
					}

				}
			}(p)
		}
	}

}

type localVisitor struct {
	p *Pair
	s *Syncer
}

func (v *localVisitor) Visit(u uri.Uri) error {
	// fmt.Println("visit", u.Uri())
	if !u.IsDir() {
		v.s.msg <- Message{u.Abs(), v.p, fsnotify.Write}
		return nil
	}
	if filepath.Base(u.Abs()) == ".git" {
		fmt.Println("skip .git", u.Uri())
		return filepath.SkipDir
	}
	v.p.watcher.Add(u.Abs())
	v.s.msg <- Message{u.Abs(), v.p, fsnotify.Create}
	return nil
}

func New() *Syncer {
	return &Syncer{
		Pairs:          make([]*Pair, 4),
		globalHandlers: make(map[fsnotify.Op][]OpHandler, 5),
		errHandlers:    make([]interface{}, 4),
		errs:           make(chan error, 4),
		msg:            make(chan Message, 4),
		addPair:        make(chan *Pair, 4),
	}
}

func (s *Syncer) Run() {
	go s.BeginWatch()
	go s.loopMsg()
	s.loopErr()
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
