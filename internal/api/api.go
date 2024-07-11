package api

import (
	"context"
	"github.com/loveuer/nf"
	"time"
	"ultone/internal/handler"
	"ultone/internal/middleware/auth"
	"ultone/internal/middleware/front"
	"ultone/internal/middleware/logger"
	"ultone/internal/middleware/oplog"
	"ultone/internal/middleware/privilege"
	"ultone/internal/model"
	"ultone/internal/opt"
)

func initApp(ctx context.Context) *nf.App {
	engine := nf.New(nf.Config{DisableLogger: true})
	engine.Use(logger.New())

	// todo: add project prefix, if you need
	// for example: app := engine.Group("/api/{project}")
	app := engine.Group("/api")
	app.Get("/available", func(c *nf.Ctx) error {
		return c.JSON(nf.Map{"status": 200, "ok": true, "time": time.Now()})
	})

	{
		api := app.Group("/user")
		api.Post("/auth/login", oplog.NewOpLog(ctx), handler.AuthLogin)
		api.Get("/auth/login", auth.NewAuth(), handler.AuthVerify)
		api.Post("/auth/logout", auth.NewAuth(), oplog.NewOpLog(ctx), handler.AuthLogout)

		mng := api.Group("/manage")
		mng.Use(auth.NewAuth(), privilege.Verify(
			privilege.RelationAnd,
			model.PrivilegeUserManage,
		))

		mng.Get("/user/list", handler.ManageUserList)
		mng.Post("/user/create", oplog.NewOpLog(ctx), handler.ManageUserCreate)
		mng.Post("/user/update", oplog.NewOpLog(ctx), handler.ManageUserUpdate)
		mng.Post("/user/delete", oplog.NewOpLog(ctx), handler.ManageUserDelete)
	}

	{
		api := app.Group("/log")
		api.Use(auth.NewAuth(), privilege.Verify(privilege.RelationAnd, model.PrivilegeOpLog))
		api.Get("/category/list", handler.LogCategories())
		api.Get("/content/list", handler.LogList)
	}

	{
		// todo: 替换 xxx
		// todo: 这里写你的模块和接口
		api := app.Group("/xxx")
		api.Use(auth.NewAuth())
		_ = api // todo: 添加自己的接口后删除该行
	}

	if opt.EnableFront {
		engine.Use(front.NewFront(&front.DefaultFront, "dist/front/browser"))
	}

	return engine
}
