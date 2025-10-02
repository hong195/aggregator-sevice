package v1

import (
	"context"
	"fmt"
	"github.com/hong195/aggregator-sevice/internal/usecase/query"

	v1 "github.com/hong195/aggregator-sevice/docs/proto/v1"
	"github.com/hong195/aggregator-sevice/internal/controller/grpc/v1/response"
)

func (r *V1) FindPacketByID(ctx context.Context, v1Request *v1.FindPacketByIDRequest) (*v1.FindPacketByIDResponse, error) {
	packet, err := r.u.Queries.FindDataPacketById.Handle(ctx, v1Request.Id)

	if err != nil {
		r.l.Error(err, "grpc - v1 - FindById - r.u.Queries.FindDataPacketById.Handle")

		return nil, fmt.Errorf("grpc - v1 - FindById %s: %w", v1Request.Id, err)
	}

	return response.NewFindByIdResponse(packet), nil
}

func (r *V1) ListPacketsByPeriod(ctx context.Context, v1Request *v1.ListPacketsByPeriodRequest) (*v1.ListPacketsByPeriodResponse, error) {
	query := query.FindDataPacketByPeriodQuery{
		Start: v1Request.Start.AsTime().UnixMilli(),
		End:   v1Request.End.AsTime().UnixMilli(),
	}
	packets, err := r.u.Queries.FindDataPacketByPeriod.Handle(ctx, query)

	if err != nil {
		r.l.Error(err, "grpc - v1 - ListPacketsByPeriod - r.u.Queries.ListDataPacketsByPeriod.Handle")

		return nil, fmt.Errorf("grpc - v1 - ListPacketsByPeriod %s - %s: %w", v1Request.Start.AsTime(), v1Request.End.AsTime(), err)
	}

	return response.NewListPacketsByPeriodResponse(packets), nil
}
