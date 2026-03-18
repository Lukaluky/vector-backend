package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	domain "vektor-backend/internal/domain/shipment"
	postgresrepo "vektor-backend/internal/repository/postgres"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@127.0.0.1:5434/shipments?sslmode=disable"
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	_, err = pool.Exec(ctx, `DELETE FROM shipment_events`)
	if err != nil {
		t.Fatalf("failed to cleanup shipment_events: %v", err)
	}

	_, err = pool.Exec(ctx, `DELETE FROM shipments`)
	if err != nil {
		t.Fatalf("failed to cleanup shipments: %v", err)
	}

	return pool
}

func TestShipmentRepository_CreateAndGetByReference(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := postgresrepo.NewShipmentRepository(pool)

	now := time.Now().UTC()

	sh, evt, err := domain.NewShipment(
		"REF-TEST-001",
		"Almaty",
		"Astana",
		"John Doe",
		"TRUCK-01",
		1000,
		700,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = repo.Create(context.Background(), sh, evt)
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	got, err := repo.GetByReference(context.Background(), "REF-TEST-001")
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}

	if got.Reference != "REF-TEST-001" {
		t.Fatalf("expected REF-TEST-001, got %s", got.Reference)
	}

	if got.CurrentStatus != domain.StatusPending {
		t.Fatalf("expected pending, got %s", got.CurrentStatus)
	}
}

func TestShipmentRepository_AddEventAndGetHistory(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := postgresrepo.NewShipmentRepository(pool)

	now := time.Now().UTC()

	sh, evt, err := domain.NewShipment(
		"REF-TEST-002",
		"Almaty",
		"Astana",
		"John Doe",
		"TRUCK-01",
		1000,
		700,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = repo.Create(context.Background(), sh, evt)
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	eventTime := now.Add(time.Minute)

	updated, err := repo.AddEvent(
		context.Background(),
		sh.Reference,
		domain.StatusPickedUp,
		domain.Event{
			Status:    domain.StatusPickedUp,
			CreatedAt: eventTime,
		},
	)
	if err != nil {
		t.Fatalf("unexpected add event error: %v", err)
	}

	if updated.CurrentStatus != domain.StatusPickedUp {
		t.Fatalf("expected picked_up, got %s", updated.CurrentStatus)
	}

	history, err := repo.GetHistory(context.Background(), sh.Reference)
	if err != nil {
		t.Fatalf("unexpected get history error: %v", err)
	}

	if len(history) != 2 {
		t.Fatalf("expected 2 history records, got %d", len(history))
	}
}