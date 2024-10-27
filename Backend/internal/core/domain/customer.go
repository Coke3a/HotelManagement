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
	DateOfBirth     *time.Time
	Gender          string
    CustomerTypeID 	uint64
	JoinDate        *time.Time
	Preferences     string
	LastVisitDate   *time.Time
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}