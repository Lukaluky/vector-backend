package shipment

import "time"

type Event struct {
	Status    Status
	CreatedAt time.Time
}