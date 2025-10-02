package v1

import (
	"github.com/go-playground/validator/v10"
	v1 "github.com/hong195/aggregator-sevice/docs/proto/v1"
	"github.com/hong195/aggregator-sevice/internal/usecase"
	"github.com/hong195/aggregator-sevice/pkg/logger"
)

// V1 -.
type V1 struct {
	v1.AggregationServiceServer

	u *usecase.UseCases
	l logger.Interface
	v *validator.Validate
}
