package entity

import (
	"errors"
	"fmt"
	"time"
)

var ErrShipmentNotFound = errors.New("shipment not found")

type Shipment struct {
	ID              string
	ReferenceNumber string
	Origin          string
	Destination     string
	CurrentStatus   Status
	DriverName      string
	UnitNumber      string
	ShipmentAmount  float64
	DriverRevenue   float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewShipment(
	id, referenceNumber, origin, destination, driverName, unitNumber string,
	shipmentAmount, driverRevenue float64,
) (*Shipment, error) {
	if referenceNumber == "" {
		return nil, fmt.Errorf("reference_number is required")
	}
	if origin == "" {
		return nil, fmt.Errorf("origin is required")
	}
	if destination == "" {
		return nil, fmt.Errorf("destination is required")
	}
	now := time.Now().UTC()
	return &Shipment{
		ID:              id,
		ReferenceNumber: referenceNumber,
		Origin:          origin,
		Destination:     destination,
		CurrentStatus:   StatusPending,
		DriverName:      driverName,
		UnitNumber:      unitNumber,
		ShipmentAmount:  shipmentAmount,
		DriverRevenue:   driverRevenue,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (s *Shipment) ApplyEvent(e ShipmentEvent) error {
	if e.Status == s.CurrentStatus {
		return fmt.Errorf("shipment is already in status %s", s.CurrentStatus)
	}
	if err := ValidateTransition(s.CurrentStatus, e.Status); err != nil {
		return err
	}
	s.CurrentStatus = e.Status
	s.UpdatedAt = time.Now().UTC()
	return nil
}
