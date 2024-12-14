package domain

import "time"

type BookingStatus int

const (
    BookingStatusUncheckIn BookingStatus = iota + 1
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
    CreatedAt    *time.Time
    UpdatedAt    *time.Time
    RoomID       uint64
    RoomTypeID   uint64
}
