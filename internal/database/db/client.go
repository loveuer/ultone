package db

import (
	"context"
	"gorm.io/gorm"
	"ultone/internal/opt"
	"ultone/internal/tool"
)

var (
	cli = &client{}
)

type client struct {
	cli   *gorm.DB
	ttype string
}

func Type() string {
	return cli.ttype
}

func New(ctxs ...context.Context) *gorm.DB {
	var ctx context.Context
	if len(ctxs) > 0 && ctxs[0] != nil {
		ctx = ctxs[0]
	} else {
		ctx = tool.Timeout(30)
	}

	session := cli.cli.Session(&gorm.Session{Context: ctx})

	if opt.Debug {
		session = session.Debug()
	}

	return session
}
