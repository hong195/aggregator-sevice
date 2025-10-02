package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hong195/aggregator-sevice/internal/entity"
	"github.com/hong195/aggregator-sevice/internal/repo"
	"time"
)

var InvalidUUIDError = errors.New("invalid UUID")

type StoreDataPacket struct {
	ID        string `json:"id"`
	Timestamp int    `json:"timestamp"` // expected as Unix milliseconds
	Payload   []int  `json:"payload"`
}

func NewStoreDataPacket(id string, timestamp int, payload []int) StoreDataPacket {
	return StoreDataPacket{
		ID:        id,
		Timestamp: timestamp,
		Payload:   payload,
	}
}

type StoreDataPacketHandler struct {
	repo repo.DataPacketRepository
}

// NewStoreDataPacketHandler constructs the use case handler with its dependency.
func NewStoreDataPacketHandler(repo repo.DataPacketRepository) *StoreDataPacketHandler {
	return &StoreDataPacketHandler{repo: repo}
}

func (h *StoreDataPacketHandler) Handle(ctx context.Context, cmd StoreDataPacket) error {

	uId, uErr := uuid.Parse(cmd.ID)

	if uErr != nil {
		return fmt.Errorf("%w: %s", InvalidUUIDError, cmd.ID)
	}

	ts := time.Unix(0, int64(cmd.Timestamp)*int64(time.Millisecond)).UTC()

	packet, err := entity.NewDataPacket(uId, ts, cmd.Payload)

	if err != nil {
		return fmt.Errorf("create datapacket: %w", err)
	}

	if err := h.repo.Store(ctx, packet); err != nil {
		return fmt.Errorf("store datapacket: insert failed: %w", err)
	}

	return nil
}
