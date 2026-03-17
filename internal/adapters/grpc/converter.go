package grpc

import (
	"errors"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nurtikaga/tms_proj/internal/core/entity"
	tmsmanagerpb "github.com/nurtikaga/tms_proj/tms-protos/go"
)

func shipmentToProto(s *entity.Shipment) *tmsmanagerpb.ShipmentResponse {
	return &tmsmanagerpb.ShipmentResponse{
		Id:        s.ID,
		Status:    domainStatusToProto(s.CurrentStatus),
		CreatedAt: timestamppb.New(s.CreatedAt),
	}
}

func eventToProto(e *entity.ShipmentEvent) *tmsmanagerpb.ShipmentEvent {
	return &tmsmanagerpb.ShipmentEvent{
		Id:         e.ID,
		Status:     domainStatusToProto(e.Status),
		Note:       e.Note,
		OccurredAt: timestamppb.New(e.OccurredAt),
	}
}

func domainStatusToProto(s entity.Status) tmsmanagerpb.ShipmentStatus {
	switch s {
	case entity.StatusPending:
		return tmsmanagerpb.ShipmentStatus_SHIPMENT_STATUS_PENDING
	case entity.StatusPickedUp:
		return tmsmanagerpb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP
	case entity.StatusInTransit:
		return tmsmanagerpb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT
	case entity.StatusDelivered:
		return tmsmanagerpb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED
	default:
		return tmsmanagerpb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED
	}
}

func protoStatusToDomain(s tmsmanagerpb.ShipmentStatus) (entity.Status, bool) {
	switch s {
	case tmsmanagerpb.ShipmentStatus_SHIPMENT_STATUS_PENDING:
		return entity.StatusPending, true
	case tmsmanagerpb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP:
		return entity.StatusPickedUp, true
	case tmsmanagerpb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT:
		return entity.StatusInTransit, true
	case tmsmanagerpb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED:
		return entity.StatusDelivered, true
	default:
		return "", false
	}
}

func containsError(err, target error) bool {
	return errors.Is(err, target)
}
