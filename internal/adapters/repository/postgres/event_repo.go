package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurtikaga/tms_proj/internal/core/entity"
)

type EventRepo struct {
	db *pgxpool.Pool
}

func NewEventRepo(db *pgxpool.Pool) *EventRepo {
	return &EventRepo{db: db}
}

const insertEvent = `
INSERT INTO shipment_events (id, shipment_id, status, note, occurred_at)
VALUES ($1,$2,$3,$4,$5)
`

func (r *EventRepo) Save(ctx context.Context, e *entity.ShipmentEvent) error {
	_, err := r.db.Exec(ctx, insertEvent,
		e.ID, e.ShipmentID, string(e.Status), e.Note, e.OccurredAt,
	)
	if err != nil {
		return fmt.Errorf("postgres: save event: %w", err)
	}
	return nil
}

const selectEventsByShipmentID = `
SELECT id, shipment_id, status, note, occurred_at
FROM shipment_events
WHERE shipment_id = $1
ORDER BY occurred_at ASC
`

func (r *EventRepo) FindByShipmentID(ctx context.Context, shipmentID string) ([]entity.ShipmentEvent, error) {
	rows, err := r.db.Query(ctx, selectEventsByShipmentID, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("postgres: query events: %w", err)
	}
	defer rows.Close()

	var events []entity.ShipmentEvent
	for rows.Next() {
		var e entity.ShipmentEvent
		var statusStr string
		if err = rows.Scan(&e.ID, &e.ShipmentID, &statusStr, &e.Note, &e.OccurredAt); err != nil {
			return nil, fmt.Errorf("postgres: scan event: %w", err)
		}
		e.Status = entity.Status(statusStr)
		events = append(events, e)
	}
	return events, rows.Err()
}
