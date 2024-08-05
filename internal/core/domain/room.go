package domain

import "time"

type Room struct {
	ID           uint64
	RoomNumber   string
	Type         string
	Description  string
	Status       string
	Floor        int
	Capacity     int
	DefaultPrice float64
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}
