package grpc

import (
	"time"

	domain "vektor-backend/internal/domain/shipment"
	shipmentpb "vektor-backend/proto"
)

func mapDomainShipmentToProto(sh *domain.Shipment) *shipmentpb.Shipment {
	if sh == nil {
		return nil
	}

	return &shipmentpb.Shipment{
		Reference:     sh.Reference,
		Origin:        sh.Origin,
		Destination:   sh.Destination,
		CurrentStatus: mapDomainStatusToProto(sh.CurrentStatus),
		Driver:        sh.Driver,
		Unit:          sh.Unit,
		Amount:        sh.Amount,
		DriverRevenue: sh.DriverRevenue,
	}
}

func mapDomainEventToProto(evt domain.Event) *shipmentpb.ShipmentEvent {
	return &shipmentpb.ShipmentEvent{
		Status:    mapDomainStatusToProto(evt.Status),
		CreatedAt: evt.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func mapProtoStatusToDomain(status shipmentpb.ShipmentStatus) domain.Status {
	switch status {
	case shipmentpb.ShipmentStatus_SHIPMENT_STATUS_PENDING:
		return domain.StatusPending
	case shipmentpb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP:
		return domain.StatusPickedUp
	case shipmentpb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT:
		return domain.StatusInTransit
	case shipmentpb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED:
		return domain.StatusDelivered
	case shipmentpb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED:
		return domain.StatusCancelled
	default:
		return ""
	}
}

func mapDomainStatusToProto(status domain.Status) shipmentpb.ShipmentStatus {
	switch status {
	case domain.StatusPending:
		return shipmentpb.ShipmentStatus_SHIPMENT_STATUS_PENDING
	case domain.StatusPickedUp:
		return shipmentpb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP
	case domain.StatusInTransit:
		return shipmentpb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT
	case domain.StatusDelivered:
		return shipmentpb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED
	case domain.StatusCancelled:
		return shipmentpb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED
	default:
		return shipmentpb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED
	}
}
