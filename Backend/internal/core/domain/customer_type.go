package domain

import "time"

type CustomerType struct {
    ID          uint64
    Name        string
    Description string
    CreatedAt   *time.Time
    UpdatedAt   *time.Time
}
