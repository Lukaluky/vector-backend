package shipment

import (
	"context"
	"time"

	domain "vektor-backend/internal/domain/shipment"
)

type CreateUseCase struct {
	repo domain.Repository
}

func NewCreateUseCase(repo domain.Repository) *CreateUseCase {
	return &CreateUseCase{repo: repo}
}

type CreateInput struct {
	Reference     string
	Origin        string
	Destination   string
	Driver        string
	Unit          string
	Amount        float64
	DriverRevenue float64
}

func (uc *CreateUseCase) Execute(
	ctx context.Context,
	input CreateInput,
) (*domain.Shipment, error) {

	now := time.Now().UTC()

	shipment, initialEvent, err := domain.NewShipment(
		input.Reference,
		input.Origin,
		input.Destination,
		input.Driver,
		input.Unit,
		input.Amount,
		input.DriverRevenue,
		now,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, shipment, initialEvent); err != nil {
		return nil, err
	}

	return shipment, nil
}
