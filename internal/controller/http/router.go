// Package v1 implements routing paths. Each services in own file.
package http

import (
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/hong195/aggregator-sevice/config"
	_ "github.com/hong195/aggregator-sevice/docs" // Swagger docs.
	"github.com/hong195/aggregator-sevice/internal/controller/http/middleware"
	v1 "github.com/hong195/aggregator-sevice/internal/controller/http/v1"
	"github.com/hong195/aggregator-sevice/internal/usecase"
	"github.com/hong195/aggregator-sevice/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Data aggregation service Rest API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /
func NewRouter(app *fiber.App, cfg *config.Config, t *usecase.UseCases, l logger.Interface) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	// metrics
	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("data-aggregator-service")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// health check
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// Routers
	apiV1Group := app.Group("api/v1")
	{
		v1.NewDataPacketRotes(apiV1Group, t, l)
	}
}
