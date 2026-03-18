package shipment

import "errors"

var (
	ErrShipmentAlreadyExists = errors.New("shipment already exists")
	ErrShipmentNotFound      = errors.New("shipment not found")
	ErrInvalidStatus         = errors.New("invalid shipment status")
	ErrInvalidTransition     = errors.New("invalid shipment status transition")
	ErrInvalidReference      = errors.New("reference is required")
	ErrInvalidOrigin         = errors.New("origin is required")
	ErrInvalidDestination    = errors.New("destination is required")
	ErrInvalidDriver         = errors.New("driver is required")
	ErrInvalidUnit           = errors.New("unit is required")
	ErrInvalidAmount         = errors.New("amount must be greater than zero")
	ErrInvalidDriverRevenue  = errors.New("driver revenue must be greater than or equal to zero")
)
