package domain

import "time"

type RatePrice struct {
	ID                uint64
	Name              string
	Description       string
	DiscountPercentage float64
	StartDate         *time.Time
	EndDate           *time.Time
	RoomID            uint64
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
}
