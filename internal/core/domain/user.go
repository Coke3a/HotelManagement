package domain

import "time"

type User struct {
	ID        uint64
	UserName  string
	Email     *string
	Role      *string
	Rank      *string
	HireDate  *time.Time
	LastLogin *time.Time
	Status    *string
	CreatedAt time.Time
	UpdatedAt time.Time
	Password string
}
