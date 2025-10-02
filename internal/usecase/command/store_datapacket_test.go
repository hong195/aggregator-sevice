package command

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hong195/aggregator-sevice/internal/entity"
	"github.com/hong195/aggregator-sevice/internal/repo/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// УСПЕШНЫЙ сценарий: корректный cmd → посчитали max → repo.Store вызван с правильными полями.
func TestStoreDataPacketHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDataPacketRepository(ctrl)

	id := uuid.NewString()
	tsMs := time.Now().UTC().UnixMilli()
	expectedTS := time.UnixMilli(tsMs).UTC()
	payload := []int{3, 7, 1} // ожидаемый max = 7

	// Проверяем, что в repo.Store пришёл корректный агрегат.
	mockRepo.
		EXPECT().
		Store(gomock.Any(), gomock.AssignableToTypeOf(entity.DataPacket{})).
		DoAndReturn(func(_ context.Context, p entity.DataPacket) error {
			assert.Equal(t, id, p.ID.String(), "id mismatch")
			assert.True(t, p.Timestamp.Equal(expectedTS), "ts mismatch: got=%s want=%s", p.Timestamp, expectedTS)
			assert.Equal(t, 7, p.MaxValue, "max mismatch")

			return nil
		})

	h := NewStoreDataPacketHandler(mockRepo)
	cmd := NewStoreDataPacket(id, tsMs, payload)

	if err := h.Handle(context.Background(), cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// НЕВАЛИДНЫЙ UUID: repo.Store вызываться не должен, ошибка обёрнута InvalidUUIDError.
func TestStoreDataPacketHandler_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDataPacketRepository(ctrl)

	h := NewStoreDataPacketHandler(mockRepo)
	cmd := NewStoreDataPacket("not-a-uuid", time.Now().UnixMilli(), []int{1})

	err := h.Handle(context.Background(), cmd)

	assert.Error(t, err)
	assert.ErrorIs(t, err, InvalidUUIDError)
}

// ОШИБКА РЕПОЗИТОРИЯ: пробрасывается наружу.
func TestStoreDataPacketHandler_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDataPacketRepository(ctrl)

	id := uuid.NewString()
	tsMs := time.Now().UTC().UnixMilli()
	payload := []int{10}

	mockRepo.
		EXPECT().
		Store(gomock.Any(), gomock.AssignableToTypeOf(entity.DataPacket{})).
		Return(errors.New("db down"))

	h := NewStoreDataPacketHandler(mockRepo)
	err := h.Handle(context.Background(), NewStoreDataPacket(id, tsMs, payload))

	assert.Error(t, err)
}
