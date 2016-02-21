package syncer

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
func (ctx Context) EmitErr(typ ErrType, v ...interface{}) {
	go func() {
		ctx.p.syncer.errs <- &Error{typ, v}
	}()
}
