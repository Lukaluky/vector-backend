package shipment

import (
	"context"

	domain "vektor-backend/internal/domain/shipment"
)

type GetHistoryUseCase struct {
	repo domain.Repository
}

func NewGetHistoryUseCase(repo domain.Repository) *GetHistoryUseCase {
	return &GetHistoryUseCase{repo: repo}
}

func (uc *GetHistoryUseCase) Execute(
	ctx context.Context,
	reference string,
) ([]domain.Event, error) {
	return uc.repo.GetHistory(ctx, reference)
}
