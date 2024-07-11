package cmd

import (
	"context"
	"ultone/internal/api"
	"ultone/internal/controller"
	"ultone/internal/database/cache"
	"ultone/internal/database/db"
	"ultone/internal/model"
	"ultone/internal/opt"
	"ultone/internal/tool"
)

var (
	filename string
)

func Execute(ctx context.Context) error {

	tool.Must(opt.Init(filename))
	tool.Must(db.Init())
	tool.Must(cache.Init())

	// todo: if elastic search required
	// tool.Must(es.Init())

	// todo: if nebula required
	// tool.Must(nebula.Init(ctx, opt.Cfg.Nebula))

	tool.Must(model.Init(db.New()))
	tool.Must(controller.Init())
	tool.Must(api.Start(ctx))

	<-ctx.Done()

	return nil
}
