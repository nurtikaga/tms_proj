package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/nurtikaga/tms_proj/internal/core/entity"
)

type ShipmentRepo struct {
	mu        sync.RWMutex
	shipments map[string]*entity.Shipment
}

func NewShipmentRepo() *ShipmentRepo {
	return &ShipmentRepo{
		shipments: make(map[string]*entity.Shipment),
	}
}

func (r *ShipmentRepo) Save(_ context.Context, s *entity.Shipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *s
	r.shipments[s.ID] = &cp
	return nil
}

func (r *ShipmentRepo) FindByID(_ context.Context, id string) (*entity.Shipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.shipments[id]
	if !ok {
		return nil, fmt.Errorf("memory: %w, id=%s", entity.ErrShipmentNotFound, id)
	}
	cp := *s
	return &cp, nil
}

func (r *ShipmentRepo) Update(_ context.Context, s *entity.Shipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.shipments[s.ID]; !ok {
		return fmt.Errorf("memory: %w, id=%s", entity.ErrShipmentNotFound, s.ID)
	}
	cp := *s
	r.shipments[s.ID] = &cp
	return nil
}

type EventRepo struct {
	mu     sync.RWMutex
	events map[string][]entity.ShipmentEvent
}

func NewEventRepo() *EventRepo {
	return &EventRepo{
		events: make(map[string][]entity.ShipmentEvent),
	}
}

func (r *EventRepo) Save(_ context.Context, e *entity.ShipmentEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events[e.ShipmentID] = append(r.events[e.ShipmentID], *e)
	return nil
}

func (r *EventRepo) FindByShipmentID(_ context.Context, shipmentID string) ([]entity.ShipmentEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	evts := r.events[shipmentID]
	out := make([]entity.ShipmentEvent, len(evts))
	copy(out, evts)
	return out, nil
}
