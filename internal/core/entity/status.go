package entity

import (
	"errors"
	"fmt"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusPickedUp  Status = "PICKED_UP"
	StatusInTransit Status = "IN_TRANSIT"
	StatusDelivered Status = "DELIVERED"
)

var ErrInvalidTransition = errors.New("invalid status transition")

var allowedTransitions = map[Status][]Status{
	StatusPending:   {StatusPickedUp},
	StatusPickedUp:  {StatusInTransit},
	StatusInTransit: {StatusDelivered},
	StatusDelivered: {},
}

func (s Status) CanTransitionTo(next Status) bool {
	allowed, ok := allowedTransitions[s]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == next {
			return true
		}
	}
	return false
}

func ValidateTransition(current, next Status) error {
	if !current.CanTransitionTo(next) {
		return fmt.Errorf("%w: %s → %s", ErrInvalidTransition, current, next)
	}
	return nil
}
