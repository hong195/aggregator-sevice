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

// DataPacketView is a read-model returned by the query handler.
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

func (h *FindDataPacketByIDHandler) Handle(ctx context.Context, id string) (DataPacketView, error) {
	uid, eUid := uuid.Parse(id)
	if eUid != nil {
		return DataPacketView{}, ErrInvalidID
	}

	p, err := h.repo.FindById(ctx, uid)
	if err != nil {
		return DataPacketView{}, fmt.Errorf("repo.FindByID: %w", err)
	}

	return DataPacketView{
		ID:        p.ID.String(),
		Timestamp: p.Timestamp.UTC(),
		MaxValue:  p.MaxValue,
	}, nil
}
