package domain

import "time"

type BookingStatus int

const (
    BookingStatusPending BookingStatus = iota + 1
    BookingStatusConfirmed
    BookingStatusCheckedIn
    BookingStatusCheckedOut
    BookingStatusCanceled
    BookingStatusCompleted
)

type Booking struct {
    ID           uint64
    CustomerID   uint64
    RatePriceId  uint64
    CheckInDate  *time.Time
    CheckOutDate *time.Time
    Status       BookingStatus
    TotalAmount  float64
    BookingDate  *time.Time
    CreatedAt    *time.Time
    UpdatedAt    *time.Time
    RoomID       uint64
    RoomTypeID   uint64
}
