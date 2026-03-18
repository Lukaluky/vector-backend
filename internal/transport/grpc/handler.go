package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	gogrpcstatus "google.golang.org/grpc/status"

	domain "vektor-backend/internal/domain/shipment"
	shipmentusecase "vektor-backend/internal/usecase/shipment"
	shipmentpb "vektor-backend/proto"
)

type Handler struct {
	shipmentpb.UnimplementedShipmentServiceServer

	createUseCase     *shipmentusecase.CreateUseCase
	getUseCase        *shipmentusecase.GetUseCase
	addEventUseCase   *shipmentusecase.AddEventUseCase
	getHistoryUseCase *shipmentusecase.GetHistoryUseCase
}

func NewHandler(
	createUseCase *shipmentusecase.CreateUseCase,
	getUseCase *shipmentusecase.GetUseCase,
	addEventUseCase *shipmentusecase.AddEventUseCase,
	getHistoryUseCase *shipmentusecase.GetHistoryUseCase,
) *Handler {
	return &Handler{
		createUseCase:     createUseCase,
		getUseCase:        getUseCase,
		addEventUseCase:   addEventUseCase,
		getHistoryUseCase: getHistoryUseCase,
	}
}

func (h *Handler) CreateShipment(
	ctx context.Context,
	req *shipmentpb.CreateShipmentRequest,
) (*shipmentpb.CreateShipmentResponse, error) {
	sh, err := h.createUseCase.Execute(ctx, shipmentusecase.CreateInput{
		Reference:     req.GetReference(),
		Origin:        req.GetOrigin(),
		Destination:   req.GetDestination(),
		Driver:        req.GetDriver(),
		Unit:          req.GetUnit(),
		Amount:        req.GetAmount(),
		DriverRevenue: req.GetDriverRevenue(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &shipmentpb.CreateShipmentResponse{
		Shipment: mapDomainShipmentToProto(sh),
	}, nil
}

func (h *Handler) GetShipment(
	ctx context.Context,
	req *shipmentpb.GetShipmentRequest,
) (*shipmentpb.GetShipmentResponse, error) {
	sh, err := h.getUseCase.Execute(ctx, req.GetReference())
	if err != nil {
		return nil, mapError(err)
	}

	return &shipmentpb.GetShipmentResponse{
		Shipment: mapDomainShipmentToProto(sh),
	}, nil
}

func (h *Handler) AddShipmentEvent(
	ctx context.Context,
	req *shipmentpb.AddShipmentEventRequest,
) (*shipmentpb.AddShipmentEventResponse, error) {
	status := mapProtoStatusToDomain(req.GetStatus())

	sh, err := h.addEventUseCase.Execute(ctx, shipmentusecase.AddEventInput{
		Reference: req.GetReference(),
		Status:    status,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &shipmentpb.AddShipmentEventResponse{
		Shipment: mapDomainShipmentToProto(sh),
	}, nil
}

func (h *Handler) GetShipmentHistory(
	ctx context.Context,
	req *shipmentpb.GetShipmentHistoryRequest,
) (*shipmentpb.GetShipmentHistoryResponse, error) {
	events, err := h.getHistoryUseCase.Execute(ctx, req.GetReference())
	if err != nil {
		return nil, mapError(err)
	}

	result := make([]*shipmentpb.ShipmentEvent, 0, len(events))
	for _, evt := range events {
		result = append(result, mapDomainEventToProto(evt))
	}

	return &shipmentpb.GetShipmentHistoryResponse{
		Events: result,
	}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrShipmentAlreadyExists):
		return gogrpcstatus.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, domain.ErrShipmentNotFound):
		return gogrpcstatus.Error(codes.NotFound, err.Error())

	case errors.Is(err, domain.ErrInvalidStatus),
		errors.Is(err, domain.ErrInvalidTransition),
		errors.Is(err, domain.ErrInvalidReference),
		errors.Is(err, domain.ErrInvalidOrigin),
		errors.Is(err, domain.ErrInvalidDestination),
		errors.Is(err, domain.ErrInvalidDriver),
		errors.Is(err, domain.ErrInvalidUnit),
		errors.Is(err, domain.ErrInvalidAmount),
		errors.Is(err, domain.ErrInvalidDriverRevenue):
		return gogrpcstatus.Error(codes.InvalidArgument, err.Error())

	default:
		return gogrpcstatus.Error(codes.Internal, err.Error())
	}
}
