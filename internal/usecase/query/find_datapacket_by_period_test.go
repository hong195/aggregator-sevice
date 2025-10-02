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

	// Период
	start := time.Now().Add(-30 * time.Minute).UTC()
	end := start.Add(10 * time.Minute).UTC()

	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)

	// Два пакета в этом периоде
	p1 := entity.DataPacket{ID: uuid.New(), Timestamp: start.Add(2 * time.Minute), MaxValue: 5}
	p2 := entity.DataPacket{ID: uuid.New(), Timestamp: start.Add(5 * time.Minute), MaxValue: 9}

	// Ожидание вызова FindByPeriod
	mockRepo.
		EXPECT().
		FindByPeriod(gomock.Any(), gomock.AssignableToTypeOf(repo.DataPacketCriteria{})).
		DoAndReturn(func(_ context.Context, c repo.DataPacketCriteria) ([]entity.DataPacket, error) {
			assert.Equal(t, start.Format(time.RFC3339), c.Start.Format(time.RFC3339))
			assert.Equal(t, end.Format(time.RFC3339), c.End.Format(time.RFC3339))
			return []entity.DataPacket{p1, p2}, nil
		})

	h := NewFindDataPacketByPeriodHandler(mockRepo)
	out, err := h.Handle(context.Background(), FindDataPacketByPeriodQuery{
		Start: startStr,
		End:   endStr,
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

func TestFindDataPacketByPeriodHandler_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDataPacketRepository(ctrl)

	start := time.Now().Add(-15 * time.Minute).UTC()
	end := time.Now().UTC()

	repoErr := errors.New("db down")

	mockRepo.
		EXPECT().
		FindByPeriod(gomock.Any(), gomock.AssignableToTypeOf(repo.DataPacketCriteria{})).
		DoAndReturn(func(_ context.Context, c repo.DataPacketCriteria) ([]entity.DataPacket, error) {
			assert.Equal(t, start.Format(time.RFC3339), c.Start.Format(time.RFC3339))
			assert.Equal(t, end.Format(time.RFC3339), c.End.Format(time.RFC3339))
			return nil, repoErr
		})

	h := NewFindDataPacketByPeriodHandler(mockRepo)
	out, err := h.Handle(context.Background(), FindDataPacketByPeriodQuery{
		Start: start.Format(time.RFC3339),
		End:   end.Format(time.RFC3339),
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, out)
}

func TestFindDataPacketByPeriodHandler_InvalidPeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDataPacketRepository(ctrl)

	mockRepo.
		EXPECT().
		FindByPeriod(gomock.Any(), gomock.Any()).
		Times(0)

	// здесь end раньше start
	start := time.Now().UTC()
	end := start.Add(-5 * time.Minute).UTC()

	h := NewFindDataPacketByPeriodHandler(mockRepo)
	out, err := h.Handle(context.Background(), FindDataPacketByPeriodQuery{
		Start: start.Format(time.RFC3339),
		End:   end.Format(time.RFC3339),
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, repo.ErrInvalidPeriod)
	assert.Nil(t, out)
}
