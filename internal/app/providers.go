package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"vektor-backend/internal/config"
	domain "vektor-backend/internal/domain/shipment"
	postgresrepo "vektor-backend/internal/repository/postgres"
	grpctransport "vektor-backend/internal/transport/grpc"
	shipmentusecase "vektor-backend/internal/usecase/shipment"
)

type shipmentUseCases struct {
	create     *shipmentusecase.CreateUseCase
	get        *shipmentusecase.GetUseCase
	addEvent   *shipmentusecase.AddEventUseCase
	getHistory *shipmentusecase.GetHistoryUseCase
}

func newPostgresPool(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.PostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}

func newShipmentRepository(pool *pgxpool.Pool) domain.Repository {
	return postgresrepo.NewShipmentRepository(pool)
}

func newShipmentUseCases(repo domain.Repository) shipmentUseCases {
	return shipmentUseCases{
		create:     shipmentusecase.NewCreateUseCase(repo),
		get:        shipmentusecase.NewGetUseCase(repo),
		addEvent:   shipmentusecase.NewAddEventUseCase(repo),
		getHistory: shipmentusecase.NewGetHistoryUseCase(repo),
	}
}

func newShipmentHandler(ucs shipmentUseCases) *grpctransport.Handler {
	return grpctransport.NewHandler(
		ucs.create,
		ucs.get,
		ucs.addEvent,
		ucs.getHistory,
	)
}

func newGRPCServer(cfg config.Config, handler *grpctransport.Handler) *grpctransport.Server {
	return grpctransport.NewServer(cfg.GRPCPort, handler)
}
