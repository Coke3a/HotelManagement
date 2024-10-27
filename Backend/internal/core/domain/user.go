package domain

import "time"

type UserRole int

const (
	UserRoleStaff UserRole = iota
	UserRoleAdmin
)


type User struct {
	ID        uint64
	UserName  string
	Role      int
	Rank      *string
	HireDate  *time.Time
	LastLogin *time.Time
	Status    *string
	CreatedAt time.Time
	UpdatedAt time.Time
	Password string
}
