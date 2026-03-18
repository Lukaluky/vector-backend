package shipment

import (
	"context"
	"time"

	domain "vektor-backend/internal/domain/shipment"
)

type AddEventUseCase struct {
	repo domain.Repository
}

func NewAddEventUseCase(repo domain.Repository) *AddEventUseCase {
	return &AddEventUseCase{repo: repo}
}

type AddEventInput struct {
	Reference string
	Status    domain.Status
}

func (uc *AddEventUseCase) Execute(
	ctx context.Context,
	input AddEventInput,
) (*domain.Shipment, error) {

	current, err := uc.repo.GetByReference(ctx, input.Reference)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	event, err := current.AddEvent(input.Status, now)
	if err != nil {
		return nil, err
	}

	return uc.repo.AddEvent(ctx, input.Reference, input.Status, event)
}
