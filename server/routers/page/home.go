package page

import (
    // "net/http"
	"github.com/tango-contrib/renders"
	"github.com/Felamande/filesync.v2/server/routers/base"
)

type HomeRouter struct {
	base.BaseTplRouter
}

func (r *HomeRouter) Get() {
    if r.Data == nil{
        r.Data = make(renders.T)
    }
	r.Data["title"] = "filesync.v2 dashboard "
    r.Tpl = "home.html"
    
    r.Render(r.Tpl,r.Data)
}
