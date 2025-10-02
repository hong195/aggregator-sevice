package query

import (
	"context"
	"fmt"
	"time"

	"github.com/hong195/aggregator-sevice/internal/repo"
)

type FindDataPacketByPeriodQuery struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type FindDataPacketByPeriodHandler struct {
	repo repo.DataPacketRepository
}

func NewFindDataPacketByPeriodHandler(r repo.DataPacketRepository) *FindDataPacketByPeriodHandler {
	return &FindDataPacketByPeriodHandler{repo: r}
}

func (h *FindDataPacketByPeriodHandler) Handle(ctx context.Context, q FindDataPacketByPeriodQuery) ([]DataPacketView, error) {

	start, err := time.Parse(time.RFC3339, q.Start)
	end, err := time.Parse(time.RFC3339, q.End)

	fmt.Printf("start: %s, end: %s\n", start, end)
	if err != nil {
		return nil, repo.ErrInvalidPeriod
	}

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
