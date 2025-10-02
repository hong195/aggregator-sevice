package v1

import (
	"github.com/go-playground/validator/v10"
	v1 "github.com/hong195/aggregator-sevice/docs/proto/v1"
	"github.com/hong195/aggregator-sevice/internal/usecase"
	"github.com/hong195/aggregator-sevice/pkg/logger"
	pbgrpc "google.golang.org/grpc"
)

// NewPacketRoutes -.
func NewPacketRoutes(app *pbgrpc.Server, t usecase.UseCases, l logger.Interface) {
	r := &V1{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}
	{
		v1.RegisterPacketServer(app, r)
	}
}
