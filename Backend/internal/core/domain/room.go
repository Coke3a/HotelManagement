package domain

import "time"

type RoomStatus int

const (
	RoomStatusAvailable RoomStatus = iota + 1
	RoomStatusMaintenance
)

type Room struct {
	ID          uint64
	RoomNumber  string
	TypeID      int
	Description string
	Status      RoomStatus
	Floor       int
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type RoomWithRoomType struct {
	ID          uint64
	RoomNumber  string
	TypeID      int
	TypeName    string
	Description string
	Status      RoomStatus
	Floor       int
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}