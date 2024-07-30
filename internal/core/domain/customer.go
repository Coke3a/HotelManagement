package domain

import "time"

type Customer struct {
	ID              uint64
	Name            string
	Email           string
	Phone           string
	Address         string
	DateOfBirth     *time.Time
	Gender          string
	MembershipStatus string
	JoinDate        *time.Time
	Preferences     string
	LastVisitDate   *time.Time
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
