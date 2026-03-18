package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"vektor-backend/internal/config"
	grpctransport "vektor-backend/internal/transport/grpc"
)

type App struct {
	cfg        config.Config
	dbPool     *pgxpool.Pool
	grpcServer *grpctransport.Server
}

func New() (*App, error) {
	cfg := config.Load()
	ctx := context.Background()

	dbPool, err := newPostgresPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	repo := newShipmentRepository(dbPool)
	useCases := newShipmentUseCases(repo)
	handler := newShipmentHandler(useCases)
	grpcServer := newGRPCServer(cfg, handler)

	return &App{
		cfg:        cfg,
		dbPool:     dbPool,
		grpcServer: grpcServer,
	}, nil
}
