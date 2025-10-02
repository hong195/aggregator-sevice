// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"
	"github.com/google/uuid"

	"github.com/hong195/aggregator-sevice/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	DataPacketRepository interface {
		Store(context.Context, entity.DataPacket) error
		FindById(ctx context.Context, id uuid.UUID) (entity.DataPacket, error)
		FindByPeriod(context.Context, DataPacketCriteria) ([]entity.DataPacket, error)
	}
)
