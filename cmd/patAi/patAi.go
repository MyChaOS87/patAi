package main

import (
	"github.com/MyChaOS87/patAi/internal/api/patents"
	"github.com/MyChaOS87/patAi/internal/api/server"
	"github.com/MyChaOS87/patAi/internal/authorization"
	"github.com/MyChaOS87/patAi/internal/cmd"
	"github.com/MyChaOS87/patAi/internal/simulation"
	"github.com/MyChaOS87/patAi/pkg/log"
)

func main() {
	ctx, cancel, cfg := cmd.Init()
	defer cancel()

	simulation := simulation.NewInMemoryQueueAndQuotaServiceSimulation()
	usecase := patents.NewValuationJobUseCase(simulation, simulation)
	handler := patents.NewHandler(usecase)
	patentsRouter := patents.NewPatentsRouter(authorization.NewMockProvider(), handler)

	srv := server.NewServer(
		server.API(&cfg.API),
		server.ChildRouters(patentsRouter),
	)
	if err := srv.Run(ctx); err != nil {
		log.Errorf("error running server: %v", err)
		cancel()
	}

	<-ctx.Done()

	log.Infof("context done: %s", ctx.Err().Error())
}
