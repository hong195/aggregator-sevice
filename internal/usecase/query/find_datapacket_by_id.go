package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hong195/aggregator-sevice/internal/repo"
	"time"
)

var ErrInvalidID = errors.New("find_datapacket_by_id: invalid UUID")

// FindDataPacketByIDQuery carries the input for the query use case.
type FindDataPacketByIDQuery struct {
	ID string
}

// DataPacketView is a read-model returned by the query handler.
// Note: Timestamp is UTC; MaxValue type matches domain (int).
type DataPacketView struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	MaxValue  int       `json:"max_value"`
}

type FindDataPacketByIDHandler struct {
	repo repo.DataPacketRepository
}

// NewFindDataPacketByIDHandler constructs a new query handler.
func NewFindDataPacketByIDHandler(r repo.DataPacketRepository) *FindDataPacketByIDHandler {
	return &FindDataPacketByIDHandler{repo: r}
}

func (h *FindDataPacketByIDHandler) Handle(ctx context.Context, q FindDataPacketByIDQuery) (DataPacketView, error) {
	id, eUid := uuid.Parse(q.ID)
	if eUid != nil {
		return DataPacketView{}, ErrInvalidID
	}

	p, err := h.repo.FindById(ctx, id)
	if err != nil {
		return DataPacketView{}, fmt.Errorf("repo.FindByID: %w", err)
	}

	return DataPacketView{
		ID:        p.ID.String(),
		Timestamp: p.Timestamp.UTC(),
		MaxValue:  p.MaxValue,
	}, nil
}
