package domain

import "time"

type Payment struct {
	ID            uint64
	BookingID     uint64
	Amount        float64
	PaymentMethod string
	PaymentDate   *time.Time
	Status        string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
