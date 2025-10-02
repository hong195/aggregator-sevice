package query

import (
	"context"
	"fmt"
	"time"

	"github.com/hong195/aggregator-sevice/internal/repo"
)

type FindDataPacketByPeriodQuery struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type FindDataPacketByPeriodHandler struct {
	repo repo.DataPacketRepository
}

func NewFindDataPacketByPeriodHandler(r repo.DataPacketRepository) *FindDataPacketByPeriodHandler {
	return &FindDataPacketByPeriodHandler{repo: r}
}

func (h *FindDataPacketByPeriodHandler) Handle(ctx context.Context, q FindDataPacketByPeriodQuery) ([]DataPacketView, error) {
	start := time.UnixMilli(q.Start).UTC()
	end := time.UnixMilli(q.End).UTC()
	if end.Before(start) {
		return nil, repo.ErrInvalidPeriod
	}

	criteria := repo.DataPacketCriteria{
		Start: start,
		End:   end,
	}

	items, err := h.repo.FindByPeriod(ctx, criteria)

	if err != nil {
		return nil, fmt.Errorf("repo.FindByPeriod: %w", err)
	}

	out := make([]DataPacketView, 0, len(items))

	for _, p := range items {
		out = append(out, DataPacketView{
			ID:        p.ID.String(),
			Timestamp: p.Timestamp.UTC(),
			MaxValue:  p.MaxValue,
		})
	}

	return out, nil
}
