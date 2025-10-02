// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"
	"errors"
	"github.com/google/uuid"

	"github.com/hong195/aggregator-sevice/internal/entity"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=contracts.go -destination=./mocks/repo_mock.go -package=mocks

var ErrInvalidPeriod = errors.New("invalid period")

type (
	DataPacketRepository interface {
		Store(context.Context, entity.DataPacket) error
		FindById(ctx context.Context, id uuid.UUID) (entity.DataPacket, error)
		FindByPeriod(context.Context, DataPacketCriteria) ([]entity.DataPacket, error)
	}
)
