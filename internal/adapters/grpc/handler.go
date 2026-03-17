package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	tmsmanagerpb "github.com/nurtikaga/tms_proj/tms-protos/go"

	"github.com/nurtikaga/tms_proj/internal/core/entity"
	"github.com/nurtikaga/tms_proj/internal/core/port"
)

type Handler struct {
	tmsmanagerpb.UnimplementedShipmentServiceServer
	svc port.ShipmentService
}

func NewHandler(svc port.ShipmentService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateShipment(ctx context.Context, req *tmsmanagerpb.CreateShipmentRequest) (*tmsmanagerpb.ShipmentResponse, error) {
	if req.GetReferenceNumber() == "" || req.GetOrigin() == "" || req.GetDestination() == "" {
		return nil, status.Error(codes.InvalidArgument, "reference_number, origin and destination are required")
	}

	input := port.CreateShipmentInput{
		ReferenceNumber: req.GetReferenceNumber(),
		Origin:          req.GetOrigin(),
		Destination:     req.GetDestination(),
		DriverName:      req.GetDriverName(),
		UnitNumber:      req.GetUnitNumber(),
		ShipmentAmount:  req.GetShipmentAmount(),
		DriverRevenue:   req.GetDriverRevenue(),
	}

	s, err := h.svc.CreateShipment(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create shipment: %v", err)
	}
	return shipmentToProto(s), nil
}

func (h *Handler) GetShipment(ctx context.Context, req *tmsmanagerpb.GetShipmentRequest) (*tmsmanagerpb.ShipmentResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	s, err := h.svc.GetShipment(ctx, req.GetId())
	if err != nil {
		if isNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "shipment %s not found", req.GetId())
		}
		return nil, status.Errorf(codes.Internal, "get shipment: %v", err)
	}
	return shipmentToProto(s), nil
}

func (h *Handler) AddEvent(ctx context.Context, req *tmsmanagerpb.AddEventRequest) (*emptypb.Empty, error) {
	if req.GetShipmentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "shipment_id is required")
	}
	if req.GetStatus() == tmsmanagerpb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "status is required")
	}

	domainStatus, ok := protoStatusToDomain(req.GetStatus())
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "unknown status: %v", req.GetStatus())
	}

	input := port.AddEventInput{
		ShipmentID: req.GetShipmentId(),
		Status:     domainStatus,
		Note:       req.GetNote(),
	}

	if err := h.svc.AddEvent(ctx, input); err != nil {
		if isNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "shipment %s not found", req.GetShipmentId())
		}
		if isInvalidTransition(err) {
			return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
		}
		return nil, status.Errorf(codes.Internal, "add event: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) GetHistory(ctx context.Context, req *tmsmanagerpb.GetHistoryRequest) (*tmsmanagerpb.HistoryResponse, error) {
	if req.GetShipmentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "shipment_id is required")
	}

	events, err := h.svc.GetHistory(ctx, req.GetShipmentId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get history: %v", err)
	}

	protoEvents := make([]*tmsmanagerpb.ShipmentEvent, 0, len(events))
	for i := range events {
		protoEvents = append(protoEvents, eventToProto(&events[i]))
	}
	return &tmsmanagerpb.HistoryResponse{Events: protoEvents}, nil
}

func isNotFound(err error) bool {
	return containsError(err, entity.ErrShipmentNotFound)
}

func isInvalidTransition(err error) bool {
	return containsError(err, entity.ErrInvalidTransition)
}
