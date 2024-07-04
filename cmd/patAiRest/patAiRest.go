package main

import (
	"github.com/MyChaOS87/patAi/internal/api/server"
	"github.com/MyChaOS87/patAi/internal/cmd"
	"github.com/MyChaOS87/patAi/pkg/log"
)

func main() {
	ctx, cancel, cfg := cmd.Init()
	defer cancel()

	srv := server.NewServer(
		server.API(&cfg.API),
	)
	if err := srv.Run(ctx); err != nil {
		log.Errorf("error running server: %v", err)
		cancel()
	}

	<-ctx.Done()

	log.Infof("context done: %s", ctx.Err().Error())
}
