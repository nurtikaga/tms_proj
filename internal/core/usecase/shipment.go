package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nurtikaga/tms_proj/internal/core/entity"
	"github.com/nurtikaga/tms_proj/internal/core/port"
)

const cacheTTL = 5 * time.Minute

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type shipmentUseCase struct {
	shipmentRepo port.ShipmentRepository
	eventRepo    port.EventRepository
	publisher    port.EventPublisher
	cache        port.ShipmentCache
	log          Logger
}

func New(
	shipmentRepo port.ShipmentRepository,
	eventRepo port.EventRepository,
	publisher port.EventPublisher,
	cache port.ShipmentCache,
	log Logger,
) port.ShipmentService {
	return &shipmentUseCase{
		shipmentRepo: shipmentRepo,
		eventRepo:    eventRepo,
		publisher:    publisher,
		cache:        cache,
		log:          log,
	}
}

func (uc *shipmentUseCase) CreateShipment(ctx context.Context, input port.CreateShipmentInput) (*entity.Shipment, error) {
	id := uuid.NewString()
	s, err := entity.NewShipment(
		id,
		input.ReferenceNumber,
		input.Origin,
		input.Destination,
		input.DriverName,
		input.UnitNumber,
		input.ShipmentAmount,
		input.DriverRevenue,
	)
	if err != nil {
		return nil, fmt.Errorf("build shipment: %w", err)
	}

	if err = uc.shipmentRepo.Save(ctx, s); err != nil {
		return nil, fmt.Errorf("save shipment: %w", err)
	}

	if cerr := uc.cache.Set(ctx, s, cacheTTL); cerr != nil {
		uc.log.Error("cache.Set failed after create", "error", cerr)
	}

	uc.log.Info("shipment created", "id", s.ID, "reference", s.ReferenceNumber)
	return s, nil
}

func (uc *shipmentUseCase) GetShipment(ctx context.Context, id string) (*entity.Shipment, error) {
	if cached, err := uc.cache.Get(ctx, id); err == nil && cached != nil {
		return cached, nil
	}

	s, err := uc.shipmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find shipment: %w", err)
	}

	if cerr := uc.cache.Set(ctx, s, cacheTTL); cerr != nil {
		uc.log.Error("cache.Set failed after get", "error", cerr)
	}
	return s, nil
}

func (uc *shipmentUseCase) AddEvent(ctx context.Context, input port.AddEventInput) error {
	s, err := uc.shipmentRepo.FindByID(ctx, input.ShipmentID)
	if err != nil {
		return fmt.Errorf("find shipment: %w", err)
	}

	event := entity.ShipmentEvent{
		ID:         uuid.NewString(),
		ShipmentID: input.ShipmentID,
		Status:     input.Status,
		Note:       input.Note,
		OccurredAt: time.Now().UTC(),
	}

	if err = s.ApplyEvent(event); err != nil {
		return err
	}

	if err = uc.eventRepo.Save(ctx, &event); err != nil {
		return fmt.Errorf("save event: %w", err)
	}

	if err = uc.shipmentRepo.Update(ctx, s); err != nil {
		return fmt.Errorf("update shipment: %w", err)
	}

	if cerr := uc.cache.Delete(ctx, s.ID); cerr != nil {
		uc.log.Error("cache.Delete failed after event", "error", cerr)
	}

	if perr := uc.publisher.Publish(ctx, &event); perr != nil {
		uc.log.Error("event publish failed", "error", perr)
	}

	uc.log.Info("shipment event added", "shipment_id", s.ID, "status", event.Status)
	return nil
}

func (uc *shipmentUseCase) GetHistory(ctx context.Context, shipmentID string) ([]entity.ShipmentEvent, error) {
	events, err := uc.eventRepo.FindByShipmentID(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}
	return events, nil
}
