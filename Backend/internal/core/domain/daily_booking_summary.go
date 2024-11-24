package domain

import "time"

type SummaryStatus int

const (
    SummaryStatusUnchecked SummaryStatus = iota
    SummaryStatusChecked
    SummaryStatusConfirmed
)

type DailyBookingSummary struct {
    SummaryDate     time.Time
    TotalBookings   int
    TotalAmount     float64
    PendingBookings int
    ConfirmedBookings int
    CheckedInBookings int
    CheckedOutBookings int
    CanceledBookings int
    CompletedBookings int
    BookingIDs      string // Store booking IDs as a comma-separated string
    Status          SummaryStatus
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
