package models

import "time"

type UserSex int

const (
	Male UserSex = iota
	Female
)

type Profile struct {
	ID          int64
	UserID      int64
	Avatar      string
	IsActive    bool
	Description string
	Username    string
	Name        string
	Surname     string
	Patronomyc  *string
	Email       string
	PhoneNumber *string
	Sex         *UserSex
	Profession  *string
	CaseCount   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserSocialLink struct {
	ID     int64
	UserID int64
	Type   string
	URL    string
}

type UserPurpose struct {
	ID      int64
	UserID  int64
	Purpose string
}
