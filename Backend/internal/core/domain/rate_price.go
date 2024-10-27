package domain

import "time"

type RatePrice struct {
	ID            uint64
	Name          string
	Description   string
	PricePerNight float64
	RoomTypeID    uint64
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
