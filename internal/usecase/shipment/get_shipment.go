package shipment

import (
	"context"

	domain "vektor-backend/internal/domain/shipment"
)

type GetUseCase struct {
	repo domain.Repository
}

func NewGetUseCase(repo domain.Repository) *GetUseCase {
	return &GetUseCase{repo: repo}
}

func (uc *GetUseCase) Execute(
	ctx context.Context,
	reference string,
) (*domain.Shipment, error) {
	return uc.repo.GetByReference(ctx, reference)
}
