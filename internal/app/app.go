package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func (a *App) Run() error {
	serverErr := make(chan error, 1)

	go func() {
		log.Printf("gRPC server started on port %s", a.cfg.GRPCPort)
		serverErr <- a.grpcServer.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case err := <-serverErr:
		if err != nil {
			a.dbPool.Close()
			return fmt.Errorf("server error: %w", err)
		}
		a.dbPool.Close()
		return nil

	case sig := <-stop:
		log.Printf("received shutdown signal: %s", sig)
	}

	log.Println("starting graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.grpcServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("grpc shutdown error: %v", err)
	}

	a.dbPool.Close()

	log.Println("application stopped")

	return nil
}
