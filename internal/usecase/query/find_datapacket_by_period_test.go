package query

import (
	"context"
	"errors"
	"github.com/hong195/aggregator-sevice/internal/repo/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/google/uuid"
	"github.com/hong195/aggregator-sevice/internal/entity"
	"github.com/hong195/aggregator-sevice/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestFindDataPacketByPeriodHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDataPacketRepository(ctrl)

	// Готовим период в ms и такие же time.Time из них (точность совпадёт с хэндлером)
	startMs := time.Now().Add(-30 * time.Minute).UTC().UnixMilli()
	endMs := startMs + int64(10*time.Minute/time.Millisecond)
	expectedStart := time.UnixMilli(startMs).UTC()
	expectedEnd := time.UnixMilli(endMs).UTC()

	// Два доменных пакета внутри периода
	p1 := entity.DataPacket{ID: uuid.New(), Timestamp: expectedStart.Add(2 * time.Minute), MaxValue: 5}
	p2 := entity.DataPacket{ID: uuid.New(), Timestamp: expectedStart.Add(5 * time.Minute), MaxValue: 9}

	// Проверяем, что критерии дошли корректно и возвращаем данные
	mockRepo.
		EXPECT().
		FindByPeriod(gomock.Any(), gomock.AssignableToTypeOf(repo.DataPacketCriteria{})).
		DoAndReturn(func(_ context.Context, c repo.DataPacketCriteria) ([]entity.DataPacket, error) {
			assert.True(t, c.Start.Equal(expectedStart), "start mismatch")
			assert.True(t, c.End.Equal(expectedEnd), "end mismatch")
			return []entity.DataPacket{p1, p2}, nil
		})

	h := NewFindDataPacketByPeriodHandler(mockRepo)
	out, err := h.Handle(context.Background(), FindDataPacketByPeriodQuery{
		Start: startMs,
		End:   endMs,
	})

	assert.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, p1.ID.String(), out[0].ID)
	assert.True(t, out[0].Timestamp.Equal(p1.Timestamp))
	assert.Equal(t, p1.MaxValue, out[0].MaxValue)
	assert.Equal(t, p2.ID.String(), out[1].ID)
	assert.True(t, out[1].Timestamp.Equal(p2.Timestamp))
	assert.Equal(t, p2.MaxValue, out[1].MaxValue)
}

func TestFindDataPacketByPeriodHandler_InvalidPeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDataPacketRepository(ctrl)

	mockRepo.
		EXPECT().
		FindByPeriod(gomock.Any(), gomock.Any()).
		Times(0)

	startMs := time.Now().UTC().UnixMilli()
	endMs := startMs - 1000

	h := NewFindDataPacketByPeriodHandler(mockRepo)
	out, err := h.Handle(context.Background(), FindDataPacketByPeriodQuery{
		Start: startMs,
		End:   endMs,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, repo.ErrInvalidPeriod)
	assert.Nil(t, out)
}

func TestFindDataPacketByPeriodHandler_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDataPacketRepository(ctrl)

	startMs := time.Now().Add(-15 * time.Minute).UTC().UnixMilli()
	endMs := time.Now().UTC().UnixMilli()
	expectedStart := time.UnixMilli(startMs).UTC()
	expectedEnd := time.UnixMilli(endMs).UTC()

	repoErr := errors.New("db down")

	mockRepo.
		EXPECT().
		FindByPeriod(gomock.Any(), gomock.AssignableToTypeOf(repo.DataPacketCriteria{})).
		DoAndReturn(func(_ context.Context, c repo.DataPacketCriteria) ([]entity.DataPacket, error) {
			assert.True(t, c.Start.Equal(expectedStart))
			assert.True(t, c.End.Equal(expectedEnd))
			return nil, repoErr
		})

	h := NewFindDataPacketByPeriodHandler(mockRepo)
	out, err := h.Handle(context.Background(), FindDataPacketByPeriodQuery{
		Start: startMs,
		End:   endMs,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, out)
}
