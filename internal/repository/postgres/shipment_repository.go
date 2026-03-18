package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	shipmentdomain "vektor-backend/internal/domain/shipment"
)

type ShipmentRepository struct {
	pool *pgxpool.Pool
}

func NewShipmentRepository(pool *pgxpool.Pool) *ShipmentRepository {
	return &ShipmentRepository{
		pool: pool,
	}
}

func (r *ShipmentRepository) Create(
	ctx context.Context,
	sh *shipmentdomain.Shipment,
	initialEvent shipmentdomain.Event,
) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO shipments (
			reference,
			origin,
			destination,
			current_status,
			driver,
			unit,
			amount,
			driver_revenue,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`,
		sh.Reference,
		sh.Origin,
		sh.Destination,
		string(sh.CurrentStatus),
		sh.Driver,
		sh.Unit,
		sh.Amount,
		sh.DriverRevenue,
		sh.CreatedAt,
		sh.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return shipmentdomain.ErrShipmentAlreadyExists
		}
		return fmt.Errorf("insert shipment: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO shipment_events (
			shipment_reference,
			status,
			created_at
		)
		VALUES ($1,$2,$3)
	`,
		sh.Reference,
		string(initialEvent.Status),
		initialEvent.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert initial event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *ShipmentRepository) GetByReference(
	ctx context.Context,
	reference string,
) (*shipmentdomain.Shipment, error) {
	var sh shipmentdomain.Shipment
	var status string

	err := r.pool.QueryRow(ctx, `
		SELECT
			reference,
			origin,
			destination,
			current_status,
			driver,
			unit,
			amount,
			driver_revenue,
			created_at,
			updated_at
		FROM shipments
		WHERE reference = $1
	`, reference).Scan(
		&sh.Reference,
		&sh.Origin,
		&sh.Destination,
		&status,
		&sh.Driver,
		&sh.Unit,
		&sh.Amount,
		&sh.DriverRevenue,
		&sh.CreatedAt,
		&sh.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shipmentdomain.ErrShipmentNotFound
		}
		return nil, fmt.Errorf("get shipment by reference: %w", err)
	}

	sh.CurrentStatus = shipmentdomain.Status(status)
	if !sh.CurrentStatus.IsValid() {
		return nil, fmt.Errorf("invalid shipment status in db: %s", status)
	}

	return &sh, nil
}

func (r *ShipmentRepository) AddEvent(
	ctx context.Context,
	reference string,
	status shipmentdomain.Status,
	event shipmentdomain.Event,
) (*shipmentdomain.Shipment, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	res, err := tx.Exec(ctx, `
		UPDATE shipments
		SET current_status = $1, updated_at = $2
		WHERE reference = $3
	`,
		string(status),
		event.CreatedAt,
		reference,
	)
	if err != nil {
		return nil, fmt.Errorf("update shipment status: %w", err)
	}

	if res.RowsAffected() == 0 {
		return nil, shipmentdomain.ErrShipmentNotFound
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO shipment_events (
			shipment_reference,
			status,
			created_at
		)
		VALUES ($1,$2,$3)
	`,
		reference,
		string(event.Status),
		event.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert shipment event: %w", err)
	}

	var sh shipmentdomain.Shipment
	var currentStatus string

	err = tx.QueryRow(ctx, `
		SELECT
			reference,
			origin,
			destination,
			current_status,
			driver,
			unit,
			amount,
			driver_revenue,
			created_at,
			updated_at
		FROM shipments
		WHERE reference = $1
	`, reference).Scan(
		&sh.Reference,
		&sh.Origin,
		&sh.Destination,
		&currentStatus,
		&sh.Driver,
		&sh.Unit,
		&sh.Amount,
		&sh.DriverRevenue,
		&sh.CreatedAt,
		&sh.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get updated shipment: %w", err)
	}

	sh.CurrentStatus = shipmentdomain.Status(currentStatus)
	if !sh.CurrentStatus.IsValid() {
		return nil, fmt.Errorf("invalid shipment status in db: %s", currentStatus)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &sh, nil
}

func (r *ShipmentRepository) GetHistory(
	ctx context.Context,
	reference string,
) ([]shipmentdomain.Event, error) {
	exists, err := r.shipmentExists(ctx, reference)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, shipmentdomain.ErrShipmentNotFound
	}

	rows, err := r.pool.Query(ctx, `
		SELECT
			status,
			created_at
		FROM shipment_events
		WHERE shipment_reference = $1
		ORDER BY created_at ASC, id ASC
	`, reference)
	if err != nil {
		return nil, fmt.Errorf("query shipment history: %w", err)
	}
	defer rows.Close()

	events := make([]shipmentdomain.Event, 0)
	for rows.Next() {
		var evt shipmentdomain.Event
		var status string

		if err := rows.Scan(&status, &evt.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan shipment history row: %w", err)
		}

		evt.Status = shipmentdomain.Status(status)
		if !evt.Status.IsValid() {
			return nil, fmt.Errorf("invalid event status in db: %s", status)
		}

		events = append(events, evt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate shipment history rows: %w", err)
	}

	return events, nil
}

func (r *ShipmentRepository) shipmentExists(ctx context.Context, reference string) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM shipments
			WHERE reference = $1
		)
	`, reference).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check shipment exists: %w", err)
	}

	return exists, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
