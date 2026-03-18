package shipment

import "testing"

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{"pending", StatusPending, true},
		{"picked_up", StatusPickedUp, true},
		{"in_transit", StatusInTransit, true},
		{"delivered", StatusDelivered, true},
		{"cancelled", StatusCancelled, true},
		{"invalid", Status("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from Status
		to   Status
		want bool
	}{
		{"pending to picked_up", StatusPending, StatusPickedUp, true},
		{"pending to cancelled", StatusPending, StatusCancelled, true},
		{"pending to delivered", StatusPending, StatusDelivered, false},
		{"picked_up to in_transit", StatusPickedUp, StatusInTransit, true},
		{"in_transit to delivered", StatusInTransit, StatusDelivered, true},
		{"delivered to cancelled", StatusDelivered, StatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.from.CanTransitionTo(tt.to)
			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
