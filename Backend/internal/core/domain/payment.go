package domain

import "time"

type PaymentStatus int
type PaymentMethod int

const (
	PaymentStatusPending PaymentStatus = iota + 1
	PaymentStatusCompleted
	PaymentStatusFailed
	PaymentStatusRefunded
)

const (
	PaymentMethodNotSpecified PaymentMethod = iota
	PaymentMethodCreditCard 
	PaymentMethodDebitCard
	PaymentMethodCash
	PaymentMethodBankTransfer
)

type Payment struct {
	ID            uint64
	BookingID     uint64
	Amount        float64
	PaymentMethod PaymentMethod
	PaymentDate   *time.Time
	Status        PaymentStatus
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
