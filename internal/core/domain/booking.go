package domain

import "time"

type Booking struct {
	ID            uint64
	CustomerID    uint64
	RoomID        uint64
	CheckInDate   *time.Time
	CheckOutDate  *time.Time
	Status        string
	TotalAmount   float64
	BookingDate   *time.Time
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
