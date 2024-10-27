package domain

import "time"

type RoomType struct {
    ID           uint64
    Name         string
    Description  string
    Capacity     int
    DefaultPrice float64
    CreatedAt    *time.Time
    UpdatedAt    *time.Time
}