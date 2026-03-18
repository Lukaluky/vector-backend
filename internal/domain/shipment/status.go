package shipment

type Status string

const (
	StatusPending   Status = "pending"
	StatusPickedUp  Status = "picked_up"
	StatusInTransit Status = "in_transit"
	StatusDelivered Status = "delivered"
	StatusCancelled Status = "cancelled"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusPickedUp, StatusInTransit, StatusDelivered, StatusCancelled:
		return true
	default:
		return false
	}
}

func (s Status) CanTransitionTo(next Status) bool {
	switch s {
	case StatusPending:
		return next == StatusPickedUp || next == StatusCancelled
	case StatusPickedUp:
		return next == StatusInTransit
	case StatusInTransit:
		return next == StatusDelivered
	case StatusDelivered, StatusCancelled:
		return false
	default:
		return false
	}
}
