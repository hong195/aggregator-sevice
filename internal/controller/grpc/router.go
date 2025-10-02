package grpc

import (
	v1 "github.com/hong195/aggregator-sevice/internal/controller/grpc/v1"
	"github.com/hong195/aggregator-sevice/internal/usecase"
	"github.com/hong195/aggregator-sevice/pkg/logger"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewRouter -.
func NewRouter(app *pbgrpc.Server, u *usecase.UseCases, l logger.Interface) {
	{
		v1.NewPacketRoutes(app, u, l)
	}

	reflection.Register(app)
}
