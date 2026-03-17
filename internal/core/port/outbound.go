package port

import (
	"context"
	"time"

	"github.com/nurtikaga/tms_proj/internal/core/entity"
)

type ShipmentRepository interface {
	Save(ctx context.Context, s *entity.Shipment) error
	FindByID(ctx context.Context, id string) (*entity.Shipment, error)
	Update(ctx context.Context, s *entity.Shipment) error
}

type EventRepository interface {
	Save(ctx context.Context, e *entity.ShipmentEvent) error
	FindByShipmentID(ctx context.Context, shipmentID string) ([]entity.ShipmentEvent, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, e *entity.ShipmentEvent) error
}

type ShipmentCache interface {
	Get(ctx context.Context, id string) (*entity.Shipment, error)
	Set(ctx context.Context, s *entity.Shipment, ttl time.Duration) error
	Delete(ctx context.Context, id string) error
}
