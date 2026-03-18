package shipment

import "context"

type Repository interface {
	Create(ctx context.Context, shipment *Shipment, initialEvent Event) error
	GetByReference(ctx context.Context, reference string) (*Shipment, error)
	AddEvent(ctx context.Context, reference string, status Status, event Event) (*Shipment, error)
	GetHistory(ctx context.Context, reference string) ([]Event, error)
}
