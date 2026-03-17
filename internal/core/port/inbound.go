package port

import (
	"context"

	"github.com/nurtikaga/tms_proj/internal/core/entity"
)

type CreateShipmentInput struct {
	ReferenceNumber string
	Origin          string
	Destination     string
	DriverName      string
	UnitNumber      string
	ShipmentAmount  float64
	DriverRevenue   float64
}

type AddEventInput struct {
	ShipmentID string
	Status     entity.Status
	Note       string
}

type ShipmentService interface {
	CreateShipment(ctx context.Context, input CreateShipmentInput) (*entity.Shipment, error)
	GetShipment(ctx context.Context, id string) (*entity.Shipment, error)
	AddEvent(ctx context.Context, input AddEventInput) error
	GetHistory(ctx context.Context, shipmentID string) ([]entity.ShipmentEvent, error)
}
