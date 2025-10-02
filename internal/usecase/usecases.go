package usecase

import (
	"github.com/hong195/aggregator-sevice/internal/repo"
	"github.com/hong195/aggregator-sevice/internal/usecase/command"
	"github.com/hong195/aggregator-sevice/internal/usecase/query"
)

type queries struct {
	FindDataPacketById     query.FindDataPacketByIDHandler
	FindDataPacketByPeriod query.FindDataPacketByPeriodHandler
}

type commands struct {
	StoreDataPacket command.StoreDataPacketHandler
}

type UseCases struct {
	Queries  queries
	Commands commands
}

func NewUseCases(repo repo.DataPacketRepository) *UseCases {
	return &UseCases{
		Queries: queries{
			FindDataPacketById:     *query.NewFindDataPacketByIDHandler(repo),
			FindDataPacketByPeriod: *query.NewFindDataPacketByPeriodHandler(repo),
		},
		Commands: commands{
			StoreDataPacket: *command.NewStoreDataPacketHandler(repo),
		},
	}
}
