package entity_test

import (
	"testing"

	"github.com/nurtikaga/tms_proj/internal/core/entity"
)

func TestNewShipment_Valid(t *testing.T) {
	s, err := entity.NewShipment("id-1", "REF-001", "City A", "City B", "Ivan", "TRUCK-01", 1000, 200)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.CurrentStatus != entity.StatusPending {
		t.Errorf("expected status %q, got %q", entity.StatusPending, s.CurrentStatus)
	}
	if s.ID != "id-1" {
		t.Errorf("expected id %q, got %q", "id-1", s.ID)
	}
}

func TestNewShipment_MissingReferenceNumber(t *testing.T) {
	_, err := entity.NewShipment("id-2", "", "City A", "City B", "", "", 0, 0)
	if err == nil {
		t.Fatal("expected error for empty reference_number, got nil")
	}
}

func TestNewShipment_MissingOrigin(t *testing.T) {
	_, err := entity.NewShipment("id-3", "REF-002", "", "City B", "", "", 0, 0)
	if err == nil {
		t.Fatal("expected error for empty origin, got nil")
	}
}

func TestNewShipment_MissingDestination(t *testing.T) {
	_, err := entity.NewShipment("id-4", "REF-003", "City A", "", "", "", 0, 0)
	if err == nil {
		t.Fatal("expected error for empty destination, got nil")
	}
}

func TestApplyEvent_ValidTransitions(t *testing.T) {
	transitions := []struct {
		from entity.Status
		to   entity.Status
	}{
		{entity.StatusPending, entity.StatusPickedUp},
		{entity.StatusPickedUp, entity.StatusInTransit},
		{entity.StatusInTransit, entity.StatusDelivered},
	}

	for _, tc := range transitions {
		t.Run(string(tc.from)+"→"+string(tc.to), func(t *testing.T) {
			s, _ := entity.NewShipment("id", "REF", "A", "B", "", "", 0, 0)
			s.CurrentStatus = tc.from

			eve := entity.ShipmentEvent{ID: "e1", ShipmentID: "id", Status: tc.to}
			if err := s.ApplyEvent(eve); err != nil {
				t.Errorf("expected valid transition %s→%s but got error: %v", tc.from, tc.to, err)
			}
			if s.CurrentStatus != tc.to {
				t.Errorf("expected current status %q, got %q", tc.to, s.CurrentStatus)
			}
		})
	}
}

func TestApplyEvent_InvalidTransitions(t *testing.T) {
	cases := []struct {
		from entity.Status
		to   entity.Status
	}{
		{entity.StatusPending, entity.StatusInTransit},
		{entity.StatusPending, entity.StatusDelivered},
		{entity.StatusDelivered, entity.StatusPending},
		{entity.StatusPickedUp, entity.StatusPending},
	}

	for _, tc := range cases {
		t.Run(string(tc.from)+"→"+string(tc.to), func(t *testing.T) {
			s, _ := entity.NewShipment("id", "REF", "A", "B", "", "", 0, 0)
			s.CurrentStatus = tc.from

			eve := entity.ShipmentEvent{ID: "e1", ShipmentID: "id", Status: tc.to}
			err := s.ApplyEvent(eve)
			if err == nil {
				t.Errorf("expected error for invalid transition %s→%s, got nil", tc.from, tc.to)
			}
		})
	}
}

func TestApplyEvent_SameStatus(t *testing.T) {
	s, _ := entity.NewShipment("id", "REF", "A", "B", "", "", 0, 0)
	eve := entity.ShipmentEvent{ID: "e1", ShipmentID: "id", Status: entity.StatusPending}
	if err := s.ApplyEvent(eve); err == nil {
		t.Error("expected error for same-status transition, got nil")
	}
}

func TestCanTransitionTo(t *testing.T) {
	if entity.StatusPending.CanTransitionTo(entity.StatusPickedUp) != true {
		t.Error("PENDING→PICKED_UP should be allowed")
	}
	if entity.StatusDelivered.CanTransitionTo(entity.StatusPending) != false {
		t.Error("DELIVERED→PENDING should be forbidden")
	}
}
