package main

import (
	"context"
	"github.com/loveuer/nf/nft/log"
	"os/signal"
	"syscall"
	"ultone/internal/cmd"
	"ultone/internal/tool"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	if err := cmd.Execute(ctx); err != nil {
		log.Error("cmd.Execute: err=%v", err)
	}

	log.Warn("received quit signal...(2s)")
	<-tool.Timeout(2).Done()
}
