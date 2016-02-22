package syncer

import "github.com/Felamande/filesync.v2/syncer/uri"

type Context struct {
	p *Pair
	s *Syncer
}

func (ctx Context) Pair() Pair {
	return ctx.p.clone()
}

func (ctx Context) Syncer() *Syncer {
	return ctx.s
}
func (ctx Context) EmitLog(typ ErrType, v ...interface{}) {
	go func() {
		ctx.p.syncer.errs <- &Error{typ, v}
	}()
}

func (ctx Context) AddWatch(u uri.Uri) error {
	return ctx.p.AddWatch(u)
}

func (ctx Context) RemoveWatch(u uri.Uri) error {
	return ctx.p.watcher.Remove(u.Abs())
}

func (ctx Context) Finish() error {
	return ErrReject
}
