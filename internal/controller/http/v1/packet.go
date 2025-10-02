package v1

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/aggregator-sevice/internal/repo"
	"github.com/hong195/aggregator-sevice/internal/usecase/query"
	"net/http"
)

// @Summary     Find packet
// @Description Find packet by ID
// @ID          history
// @Tags  	    packet
// @Produce     json
// @Success     200 {object} query.DataPacketView
// @Failure     500 {object} response.Error
// @Router      /api/v1/packets/:id [get]
func (r *V1) findPacket(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	packet, err := r.t.Queries.FindDataPacketById.Handle(ctx.Context(), idStr)

	if err != nil {
		r.l.Error(err, "http - v1 - packet")

		if errors.Is(err, repo.ErrInvalidPeriod) {
			return errorResponse(ctx, http.StatusNotFound, "Packet not found")
		}

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(packet)
}

// @Summary     Get packets by period
// @Description Return packets with max values for the time window [start, end] (Unix ms, UTC).
// @ID          packets-by-period
// @Tags        packets
// @Accept      json
// @Produce     json
// @Param       start query string true  "Start time (time.RFC3339 format)"
// @Param       end   query string true  "End time (time.RFC3339 format)"
// @Success     200 {array} query.DataPacketView
// @Failure     400 {object} response.Error "missing/invalid query params"
// @Failure     422 {object} response.Error "end < start"
// @Failure     500 {object} response.Error
// @Router      /api/v1/packets [get]
func (r *V1) listPackets(ctx *fiber.Ctx) error {
	startStr := ctx.Query("start")
	endStr := ctx.Query("end")

	if startStr == "" || endStr == "" {
		return errorResponse(ctx, http.StatusBadRequest, "query params 'start' and 'end' are required (unix ms)")
	}

	fmt.Println("startStr:", startStr, "endStr:", endStr)

	q := query.FindDataPacketByPeriodQuery{
		Start: startStr,
		End:   endStr,
	}

	items, handleErr := r.t.Queries.FindDataPacketByPeriod.Handle(ctx.UserContext(), q)

	if handleErr != nil {
		if errors.Is(handleErr, repo.ErrInvalidPeriod) {
			return errorResponse(ctx, http.StatusUnprocessableEntity, "end must be >= start")
		}
		r.l.Error(handleErr, "http - v1 - listPackets")
		return errorResponse(ctx, http.StatusInternalServerError, "storage error")
	}

	return ctx.Status(http.StatusOK).JSON(items)
}
