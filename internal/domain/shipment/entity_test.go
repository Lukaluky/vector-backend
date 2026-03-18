package shipment

import (
	"testing"
	"time"
)

func TestNewShipment_Success(t *testing.T) {
	now := time.Now().UTC()

	sh, evt, err := NewShipment(
		"REF-001",
		"Almaty",
		"Astana",
		"John Doe",
		"TRUCK-01",
		1000,
		700,
		now,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if sh.Reference != "REF-001" {
		t.Fatalf("expected reference REF-001, got %s", sh.Reference)
	}

	if sh.CurrentStatus != StatusPending {
		t.Fatalf("expected status pending, got %s", sh.CurrentStatus)
	}

	if evt.Status != StatusPending {
		t.Fatalf("expected initial event status pending, got %s", evt.Status)
	}
}

func TestNewShipment_InvalidReference(t *testing.T) {
	now := time.Now().UTC()

	_, _, err := NewShipment(
		"",
		"Almaty",
		"Astana",
		"John Doe",
		"TRUCK-01",
		1000,
		700,
		now,
	)

	if err != ErrInvalidReference {
		t.Fatalf("expected ErrInvalidReference, got %v", err)
	}
}

func TestShipment_AddEvent_ValidTransition(t *testing.T) {
	now := time.Now().UTC()

	sh, _, err := NewShipment(
		"REF-001",
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

	evt, err := sh.AddEvent(StatusPickedUp, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if sh.CurrentStatus != StatusPickedUp {
		t.Fatalf("expected status picked_up, got %s", sh.CurrentStatus)
	}

	if evt.Status != StatusPickedUp {
		t.Fatalf("expected event status picked_up, got %s", evt.Status)
	}
}

func TestShipment_AddEvent_InvalidTransition(t *testing.T) {
	now := time.Now().UTC()

	sh, _, err := NewShipment(
		"REF-001",
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

	_, err = sh.AddEvent(StatusDelivered, now.Add(time.Minute))
	if err != ErrInvalidTransition {
		t.Fatalf("expected ErrInvalidTransition, got %v", err)
	}
}
