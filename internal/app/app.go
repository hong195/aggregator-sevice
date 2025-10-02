// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/hong195/aggregator-sevice/config"
	"github.com/hong195/aggregator-sevice/internal/controller/grpc"
	"github.com/hong195/aggregator-sevice/internal/controller/http"
	"github.com/hong195/aggregator-sevice/internal/repo/persistent"
	"github.com/hong195/aggregator-sevice/internal/usecase"
	"github.com/hong195/aggregator-sevice/pkg/generator"
	"github.com/hong195/aggregator-sevice/pkg/grpcserver"
	"github.com/hong195/aggregator-sevice/pkg/httpserver"
	"github.com/hong195/aggregator-sevice/pkg/logger"
	"github.com/hong195/aggregator-sevice/pkg/postgres"
	"github.com/hong195/aggregator-sevice/pkg/worker_pool"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	repo := persistent.NewDataPacketRepository(pg)
	useCases := usecase.NewUseCases(repo)

	// gRPC Server
	grpcServer := grpcserver.New(grpcserver.Port(cfg.GRPC.Port))
	grpc.NewRouter(grpcServer.App, nil, l)

	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, useCases, l)

	//Raw packet generator
	outRawPackets := make(chan generator.RawPacket)
	gen := generator.NewGenerator(time.Millisecond*100, 10, outRawPackets, l)

	//Worker pool
	wp := worker_pool.NewPool(10, outRawPackets, useCases, l)

	// Start servers
	grpcServer.Start()
	httpServer.Start()
	gen.Start()
	wp.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	err = gen.Stop()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - generator.Shutdown: %w", err))
	}

	err = wp.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - workder_pool.Shutdown: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = grpcServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}
}
