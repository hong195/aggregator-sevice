package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/aggregator-sevice/internal/usecase"
	"github.com/hong195/aggregator-sevice/pkg/logger"
)

// NewDataPacketRotes -.
func NewDataPacketRotes(apiV1Group fiber.Router, t *usecase.UseCases, l logger.Interface) {
	r := &V1{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	translationGroup := apiV1Group.Group("/packets")

	{
		translationGroup.Post("/", r.listPackets)
		translationGroup.Get("/:id", r.findPacket)
	}
}
