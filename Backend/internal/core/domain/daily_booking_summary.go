package domain

import (
	"time"
	"fmt"
	"strings"
	"strconv"
)

type SummaryStatus int

const (
    SummaryStatusUnchecked SummaryStatus = iota
    SummaryStatusChecked
    SummaryStatusConfirmed
)

type DailyBookingSummary struct {
    SummaryDate       time.Time
    CreatedBookings   string // Store booking IDs as a comma-separated string
    CompletedBookings string // Store booking IDs as a comma-separated string
    CanceledBookings  string // Store booking IDs as a comma-separated string
    TotalAmount       float64
    Status            SummaryStatus
    CreatedAt         time.Time
    UpdatedAt         time.Time
}

// Helper functions for booking IDs formatting
func FormatBookingIDs(ids []uint64) string {
    if len(ids) == 0 {
        return ""
    }
    var result string
    for i, id := range ids {
        if i > 0 {
            result += ";"
        }
        result += fmt.Sprintf("%d", id)
    }
    return result
}

func ParseBookingIDs(idsStr string) []uint64 {
    if idsStr == "" {
        return []uint64{}
    }
    
    strIDs := strings.Split(idsStr, ";")
    ids := make([]uint64, 0, len(strIDs))
    
    for _, str := range strIDs {
        if id, err := strconv.ParseUint(str, 10, 64); err == nil {
            ids = append(ids, id)
        }
    }
    return ids
}