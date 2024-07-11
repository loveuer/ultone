package nebula

import (
	"context"
	"github.com/loveuer/ngorm/v2"
	"strings"
	"ultone/internal/opt"
)

var (
	client *ngorm.Client
)

func Init(ctx context.Context, cfg opt.Nebula) error {
	var (
		err error
	)

	if client, err = ngorm.NewClient(ctx, &ngorm.Config{
		Endpoints:    strings.Split(cfg.Uri, ","),
		Username:     cfg.Username,
		Password:     cfg.Password,
		DefaultSpace: cfg.Space,
		Logger:       nil,
	}); err != nil {
		return err
	}

	return nil
}

func New(ctx context.Context, cfgs ...*ngorm.SessCfg) *ngorm.Session {
	return client.Session(cfgs...)
}
