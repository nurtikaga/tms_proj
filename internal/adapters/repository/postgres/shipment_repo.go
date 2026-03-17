package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurtikaga/tms_proj/internal/core/entity"
)

type ShipmentRepo struct {
	db *pgxpool.Pool
}

func NewShipmentRepo(db *pgxpool.Pool) *ShipmentRepo {
	return &ShipmentRepo{db: db}
}

const insertShipment = `
INSERT INTO shipments
    (id, reference_number, origin, destination, current_status,
     driver_name, unit_number, shipment_amount, driver_revenue,
     created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`

func (r *ShipmentRepo) Save(ctx context.Context, s *entity.Shipment) error {
	_, err := r.db.Exec(ctx, insertShipment,
		s.ID, s.ReferenceNumber, s.Origin, s.Destination, string(s.CurrentStatus),
		s.DriverName, s.UnitNumber, s.ShipmentAmount, s.DriverRevenue,
		s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("postgres: save shipment: %w", err)
	}
	return nil
}

const selectShipmentByID = `
SELECT id, reference_number, origin, destination, current_status,
       driver_name, unit_number, shipment_amount, driver_revenue,
       created_at, updated_at
FROM shipments WHERE id = $1
`

func (r *ShipmentRepo) FindByID(ctx context.Context, id string) (*entity.Shipment, error) {
	row := r.db.QueryRow(ctx, selectShipmentByID, id)
	s := &entity.Shipment{}
	var statusStr string
	err := row.Scan(
		&s.ID, &s.ReferenceNumber, &s.Origin, &s.Destination, &statusStr,
		&s.DriverName, &s.UnitNumber, &s.ShipmentAmount, &s.DriverRevenue,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("postgres: %w, id=%s", entity.ErrShipmentNotFound, id)
		}
		return nil, fmt.Errorf("postgres: find shipment: %w", err)
	}
	s.CurrentStatus = entity.Status(statusStr)
	return s, nil
}

const updateShipment = `
UPDATE shipments SET current_status=$2, updated_at=$3 WHERE id=$1
`

func (r *ShipmentRepo) Update(ctx context.Context, s *entity.Shipment) error {
	_, err := r.db.Exec(ctx, updateShipment, s.ID, string(s.CurrentStatus), s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("postgres: update shipment: %w", err)
	}
	return nil
}
