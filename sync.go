package syncer

import (
	"github.com/Felamande/filesync.v2/uri"
	fsnotify "gopkg.in/fsnotify.v1"
)

type Handler func(lUri uri.Uri, rUri uri.Uri) error

type PairConfig struct {
}

type Pair struct {
	Left  uri.Uri
	Right uri.Uri

	watcher  *fsnotify.Watcher
	progress chan int64
	syncer   *Syncer
	handlers map[fsnotify.Op][]Handler

	Config *PairConfig
}

func (p *Pair) Syncer() *Syncer {
	return p.syncer
}

type Message struct {
	Name string
	P    *Pair
	Op   fsnotify.Op
}

type Syncer struct {
	Pairs          []*Pair
	globalHandlers map[fsnotify.Op][]Handler
	errHandlers    []func(error)
	errs           chan error
	msg            chan Message
}

func (s *Syncer) handleOp() {
	for {
		select {
		case msg := <-s.msg:
			go func(m Message) {
				if handlers, exist := s.globalHandlers[m.Op]; exist {
					l, r, err := prepair(m.Name, m.P)
					if err != nil {
						go func() { s.errs <- err }()
						return
					}
					for _, h := range handlers {
						if h != nil {
							h(l, r)
						}
					}
				}

			}(msg)
		}
	}
}

func (s *Syncer) handleErr() {
	for {
		select {
		case err := <-s.errs:
			go func(e error) {
				for _, h := range s.errHandlers {
					h(e)
				}
			}(err)
		}
	}
}

func (s *Syncer) HandleOp(op fsnotify.Op, h ...Handler) {
	s.globalHandlers[op] = append(s.globalHandlers[op], h...)
}

func (s *Syncer) HandleError(h ...func(error)) {
	s.errHandlers = append(s.errHandlers, h...)
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
		handlers: make(map[fsnotify.Op][]Handler),
		syncer:   s,
		Config:   config,
	}
	s.Pairs = append(s.Pairs, p)

	return nil
}

func (s *Syncer) BeginWatch() {
	for _, p := range s.Pairs {
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

type localVisitor struct {
	p *Pair
	s *Syncer
}

func (v *localVisitor) Visit(u uri.Uri) error {
	if !u.IsDir() {
		return nil
	}
	err := v.p.watcher.Add(u.Abs())
	if err != nil {
		go func() { v.s.errs <- err }()
	}
}

func New() *Syncer {
	return &Syncer{
		Pairs:          make([]*Pair, 4),
		globalHandlers: make(map[fsnotify.Op][]Handler, 5),
		errHandlers:    make([]func(error), 4),
		errs:           make(chan error, 4),
		msg:            make(chan Message, 4),
	}
}

func prepair(name string, p *Pair) (l uri.Uri, r uri.Uri, err error) {

}
