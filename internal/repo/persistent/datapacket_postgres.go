package persistent

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hong195/aggregator-sevice/internal/repo"
	"github.com/jackc/pgx/v5"

	"github.com/hong195/aggregator-sevice/internal/entity"
	"github.com/hong195/aggregator-sevice/pkg/postgres"
)

const _defaultEntityCap = 64

// DataPacketRepository -.
type DataPacketRepository struct {
	*postgres.Postgres
}

// NewDataPacketRepository -.
func NewDataPacketRepository(pg *postgres.Postgres) *DataPacketRepository {
	return &DataPacketRepository{pg}
}

// FindById -.
func (r *DataPacketRepository) FindById(ctx context.Context, id uuid.UUID) (entity.DataPacket, error) {
	sql := `SELECT id, ts, max_value FROM data_packets WHERE id = $1`

	var p entity.DataPacket
	if err := r.Pool.QueryRow(ctx, sql, id).
		Scan(&p.ID, &p.Timestamp, &p.MaxValue); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return entity.DataPacket{}, fmt.Errorf("packet %s not found: %w", id, err)
		}
		return entity.DataPacket{}, fmt.Errorf("FindByID query: %w", err)
	}
	return entity.DataPacket{}, nil
}

func (r *DataPacketRepository) FindByPeriod(ctx context.Context, criteria repo.DataPacketCriteria) ([]entity.DataPacket, error) {
	start := criteria.Start.UTC()
	end := criteria.End.UTC()

	if end.Before(start) {
		return nil, fmt.Errorf("FindByPeriod: end before start")
	}

	q := `
        SELECT id, ts, max_value
        FROM data_packets
        WHERE ts >= $1 AND ts <= $2
        ORDER BY ts ASC, id ASC
    `

	rows, err := r.Pool.Query(ctx, q, start, end)
	if err != nil {
		return nil, fmt.Errorf("FindByPeriod query: %w", err)
	}
	defer rows.Close()

	out := make([]entity.DataPacket, 0, _defaultEntityCap) // небольшая предвыделенная ёмкость

	for rows.Next() {
		var p entity.DataPacket
		if err := rows.Scan(&p.ID, &p.Timestamp, &p.MaxValue); err != nil {
			return nil, fmt.Errorf("FindByPeriod scan: %w", err)
		}
		out = append(out, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindByPeriod rows: %w", err)
	}

	return out, nil
}

// Store -.
func (r *DataPacketRepository) Store(ctx context.Context, dp entity.DataPacket) error {
	sql, args, err := r.Builder.
		Insert("data_packets").
		Columns("id, timestamp, max_value").
		Values(dp.ID, dp.Timestamp.UTC(), dp.MaxValue).
		ToSql()

	if err != nil {
		return fmt.Errorf("DataPacketRepository - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DataPacketRepository - Store - r.Pool.Exec: %w", err)
	}

	return nil
}
