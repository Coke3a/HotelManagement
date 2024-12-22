package domain

import "time"

type Customer struct {
	ID              uint64
	FirstName       string
	Surname         string
	IdentityNumber  string
	Email           string
	Phone           string
	Address         string
	Gender          string
	CustomerTypeID 	uint64
	Preferences     string
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}