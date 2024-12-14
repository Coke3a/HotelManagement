package domain

import "time"

type BookingCustomerPayment struct {
	BookingID         uint64
	CustomerID        uint64
	BookingPrice      float64
	BookingStatus     BookingStatus
	CheckInDate       *time.Time
	CheckOutDate      *time.Time
	BookingCreatedAt  *time.Time
	BookingUpdatedAt  *time.Time
	RoomID            uint64
	RoomNumber        string
	RatePriceID       uint64
	Floor             uint64
	RoomTypeID        uint64
	RoomTypeName      string
	CustomerFirstName string
	CustomerSurname   string
	CustomerIdentityNumber string
	CustomerAddress string
	PaymentID         *uint64
	PaymentStatus     *uint64
	PaymentUpdateDate *time.Time
}
