package server

import (
	"github.com/Felamande/filesync.v2/server/modules/utils"
	"github.com/Felamande/filesync.v2/settings"
	"github.com/lunny/tango"

	//middlewares
	"github.com/Felamande/filesync.v2/server/middlewares/time"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/renders"

	//routers
	"github.com/Felamande/filesync.v2/server/routers/page"
	"github.com/Felamande/filesync.v2/server/routers/pairs"
)

var t *tango.Tango

func Run() {

	t = tango.New()

	t.Use(new(time.TimeHandler))
	t.Use(tango.Static(tango.StaticOptions{
		RootPath: utils.Abs(settings.Static.LocalRoot),
	}))
	t.Use(binding.Bind())
	t.Use(tango.ClassicHandlers...)
	t.Use(renders.New(renders.Options{
		Reload:      settings.Template.Reload,
		Directory:   utils.Abs(settings.Template.Home),
		Charset:     settings.Template.Charset,
		DelimsLeft:  settings.Template.DelimesLeft,
		DelimsRight: settings.Template.DelimesRight,
		Funcs:       utils.DefaultFuncs(),
	}))
	t.Use(events.Events())
	t.Group("/pair", func(g *tango.Group) {
		g.Get("/all", new(pairs.GetAllRouter))
		g.Post("/new", new(pairs.NewPairRouter))
	})
	t.Get("/", new(page.HomeRouter))

	t.Run(settings.Server.Port)
}
