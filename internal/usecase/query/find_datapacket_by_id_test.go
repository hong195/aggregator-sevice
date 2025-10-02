package query

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

func TestFindDataPacketByIDHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDataPacketRepository(ctrl)

	id := uuid.New()
	ts := time.Now().UTC()
	dp := entity.DataPacket{
		ID:        id,
		Timestamp: ts,
		MaxValue:  123,
	}

	mockRepo.
		EXPECT().
		FindById(gomock.Any(), id).
		Return(dp, nil)

	h := NewFindDataPacketByIDHandler(mockRepo)
	view, err := h.Handle(context.Background(), FindDataPacketByIDQuery{ID: id.String()})

	assert.NoError(t, err)
	assert.Equal(t, id.String(), view.ID)
	// time.Time лучше сравнивать через Equal (игнорирует монотонную часть)
	assert.True(t, view.Timestamp.Equal(ts), "timestamp mismatch: got=%s want=%s", view.Timestamp, ts)
	assert.Equal(t, 123, view.MaxValue)
}

func TestFindDataPacketByIDHandler_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDataPacketRepository(ctrl)

	h := NewFindDataPacketByIDHandler(mockRepo)
	view, err := h.Handle(context.Background(), FindDataPacketByIDQuery{ID: "not-a-uuid"})

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidID)
	assert.Equal(t, DataPacketView{}, view)
}

func TestFindDataPacketByIDHandler_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDataPacketRepository(ctrl)

	id := uuid.New()
	repoErr := errors.New("db down")

	mockRepo.
		EXPECT().
		FindById(gomock.Any(), id).
		Return(entity.DataPacket{}, repoErr)

	h := NewFindDataPacketByIDHandler(mockRepo)
	view, err := h.Handle(context.Background(), FindDataPacketByIDQuery{ID: id.String()})

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Equal(t, DataPacketView{}, view)
}
