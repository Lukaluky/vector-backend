package shipment

import "time"

type Shipment struct {
	Reference     string
	Origin        string
	Destination   string
	CurrentStatus Status
	Driver        string
	Unit          string
	Amount        float64
	DriverRevenue float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewShipment(
	reference string,
	origin string,
	destination string,
	driver string,
	unit string,
	amount float64,
	driverRevenue float64,
	now time.Time,
) (*Shipment, Event, error) {
	if reference == "" {
		return nil, Event{}, ErrInvalidReference
	}
	if origin == "" {
		return nil, Event{}, ErrInvalidOrigin
	}
	if destination == "" {
		return nil, Event{}, ErrInvalidDestination
	}
	if driver == "" {
		return nil, Event{}, ErrInvalidDriver
	}
	if unit == "" {
		return nil, Event{}, ErrInvalidUnit
	}
	if amount <= 0 {
		return nil, Event{}, ErrInvalidAmount
	}
	if driverRevenue < 0 {
		return nil, Event{}, ErrInvalidDriverRevenue
	}

	s := &Shipment{
		Reference:     reference,
		Origin:        origin,
		Destination:   destination,
		CurrentStatus: StatusPending,
		Driver:        driver,
		Unit:          unit,
		Amount:        amount,
		DriverRevenue: driverRevenue,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	initialEvent := Event{
		Status:    StatusPending,
		CreatedAt: now,
	}

	return s, initialEvent, nil
}

func (s *Shipment) AddEvent(next Status, now time.Time) (Event, error) {
	if !next.IsValid() {
		return Event{}, ErrInvalidStatus
	}

	if s.CurrentStatus == next {
		return Event{}, ErrInvalidTransition
	}

	if !s.CurrentStatus.CanTransitionTo(next) {
		return Event{}, ErrInvalidTransition
	}

	s.CurrentStatus = next
	s.UpdatedAt = now

	event := Event{
		Status:    next,
		CreatedAt: now,
	}

	return event, nil
}