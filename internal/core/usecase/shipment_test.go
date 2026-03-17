package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/nurtikaga/tms_proj/internal/adapters/repository/memory"
	"github.com/nurtikaga/tms_proj/internal/core/entity"
	"github.com/nurtikaga/tms_proj/internal/core/port"
	"github.com/nurtikaga/tms_proj/internal/core/usecase"
)

type noopPublisher struct{}

func (noopPublisher) Publish(_ context.Context, _ *entity.ShipmentEvent) error { return nil }

type noopCache struct{}

func (noopCache) Get(_ context.Context, _ string) (*entity.Shipment, error)      { return nil, nil }
func (noopCache) Set(_ context.Context, _ *entity.Shipment, _ interface{}) error { return nil }
func (noopCache) Delete(_ context.Context, _ string) error                       { return nil }

type noopLog struct{}

func (noopLog) Info(_ string, _ ...any)  {}
func (noopLog) Error(_ string, _ ...any) {}

func buildSvc(t *testing.T) (port.ShipmentService, *memory.ShipmentRepo, *memory.EventRepo) {
	t.Helper()
	shipRepo := memory.NewShipmentRepo()
	eventRepo := memory.NewEventRepo()
	svc := usecase.New(shipRepo, eventRepo, noopPublisher{}, noopCacheAdapter{}, noopLog{})
	return svc, shipRepo, eventRepo
}

type noopCacheAdapter struct{}

func (noopCacheAdapter) Get(_ context.Context, _ string) (*entity.Shipment, error) { return nil, nil }
func (noopCacheAdapter) Set(_ context.Context, _ *entity.Shipment, _ time.Duration) error {
	return nil
}
func (noopCacheAdapter) Delete(_ context.Context, _ string) error { return nil }

func TestCreateShipment(t *testing.T) {
	svc, _, _ := buildSvc(t)
	input := port.CreateShipmentInput{
		ReferenceNumber: "REF-100",
		Origin:          "Almaty",
		Destination:     "Astana",
		DriverName:      "Nurlan",
		UnitNumber:      "KZ-001",
		ShipmentAmount:  5000,
		DriverRevenue:   500,
	}
	s, err := svc.CreateShipment(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateShipment returned error: %v", err)
	}
	if s.ID == "" {
		t.Error("expected non-empty shipment ID")
	}
	if s.CurrentStatus != entity.StatusPending {
		t.Errorf("expected status PENDING, got %q", s.CurrentStatus)
	}
	if s.ReferenceNumber != "REF-100" {
		t.Errorf("expected ReferenceNumber REF-100, got %q", s.ReferenceNumber)
	}
}

func TestCreateShipment_MissingFields(t *testing.T) {
	svc, _, _ := buildSvc(t)
	_, err := svc.CreateShipment(context.Background(), port.CreateShipmentInput{})
	if err == nil {
		t.Error("expected error for empty CreateShipmentInput, got nil")
	}
}

func TestAddEvent_ValidTransition(t *testing.T) {
	svc, _, _ := buildSvc(t)
	s, _ := svc.CreateShipment(context.Background(), port.CreateShipmentInput{
		ReferenceNumber: "REF-200", Origin: "A", Destination: "B",
	})

	err := svc.AddEvent(context.Background(), port.AddEventInput{
		ShipmentID: s.ID,
		Status:     entity.StatusPickedUp,
		Note:       "picked up at warehouse",
	})
	if err != nil {
		t.Fatalf("AddEvent returned unexpected error: %v", err)
	}

	updated, err := svc.GetShipment(context.Background(), s.ID)
	if err != nil {
		t.Fatalf("GetShipment returned error: %v", err)
	}
	if updated.CurrentStatus != entity.StatusPickedUp {
		t.Errorf("expected PICKED_UP, got %q", updated.CurrentStatus)
	}
}

func TestAddEvent_InvalidTransition(t *testing.T) {
	svc, _, _ := buildSvc(t)
	s, _ := svc.CreateShipment(context.Background(), port.CreateShipmentInput{
		ReferenceNumber: "REF-300", Origin: "A", Destination: "B",
	})

	err := svc.AddEvent(context.Background(), port.AddEventInput{
		ShipmentID: s.ID,
		Status:     entity.StatusDelivered,
		Note:       "invalid jump",
	})
	if err == nil {
		t.Error("expected error for invalid transition CREATED→DELIVERED, got nil")
	}
}

func TestGetHistory(t *testing.T) {
	svc, _, _ := buildSvc(t)
	s, _ := svc.CreateShipment(context.Background(), port.CreateShipmentInput{
		ReferenceNumber: "REF-400", Origin: "A", Destination: "B",
	})

	_ = svc.AddEvent(context.Background(), port.AddEventInput{ShipmentID: s.ID, Status: entity.StatusPickedUp})
	_ = svc.AddEvent(context.Background(), port.AddEventInput{ShipmentID: s.ID, Status: entity.StatusInTransit})

	events, err := svc.GetHistory(context.Background(), s.ID)
	if err != nil {
		t.Fatalf("GetHistory returned error: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}
