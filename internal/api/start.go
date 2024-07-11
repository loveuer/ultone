package api

import (
	"context"
	"fmt"
	"github.com/loveuer/nf/nft/log"
	"net"
	"ultone/internal/opt"
	"ultone/internal/tool"
)

func Start(ctx context.Context) error {

	app := initApp(ctx)
	ready := make(chan bool)

	ln, err := net.Listen("tcp", opt.Cfg.Listen.Http)
	if err != nil {
		return fmt.Errorf("api.MustStart: net listen tcp address=%v err=%v", opt.Cfg.Listen.Http, err)
	}

	go func() {
		ready <- true

		if err = app.RunListener(ln); err != nil {
			log.Panic("api.MustStart: app run err=%v", err)
		}
	}()

	<-ready

	go func() {
		ready <- true
		<-ctx.Done()
		if err = app.Shutdown(tool.Timeout(1)); err != nil {
			log.Error("api.MustStart: app shutdown err=%v", err)
		}
	}()

	<-ready

	return nil
}
