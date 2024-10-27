package domain

import "time"

type Log struct {
	ID        uint64
	TableName string
	RecordID  uint64
	Action    string
	UserID    uint64
	CreatedAt time.Time
}