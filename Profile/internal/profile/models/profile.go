package models

import (
	"time"
)

type UserSex int

const (
	Male UserSex = iota
	Female
)

type UserRole int

const (
	Admin UserRole = iota
	User
	Creator
)

type Profile struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Avatar      string    `json:"avatar" db:"avatar"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Description string    `json:"description" db:"description"`
	Username    string    `json:"username" db:"username"`
	Name        string    `json:"name" db:"name"`
	Surname     string    `json:"surname" db:"surname"`
	Patronymic  *string   `json:"patronymic" db:"patronymic"`
	City        *string   `json:"city" db:"city"`
	Age         *int      `json:"age" db:"age"`
	Sex         *UserSex  `json:"sex" db:"sex"`
	Profession  *string   `json:"profession" db:"profession"`
	CaseCount   int       `json:"case_count" db:"case_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserSocialLink struct {
	ID     int64  `json:"id" db:"id"`
	UserID int64  `json:"user_id" db:"user_id"`
	Type   string `json:"type" db:"type"`
	URL    string `json:"url" db:"url"`
}

type UserPurpose struct {
	ID      int64  `json:"id" db:"id"`
	UserID  int64  `json:"user_id" db:"user_id"`
	Purpose string `json:"purpose" db:"purpose"`
}

type UserProfession struct {
	ID           int64 `json:"id" db:"id"`
	UserID       int64 `json:"user_id" db:"user_id"`
	ProfessionID int16 `json:"profession_id" db:"profession_id"`
}

type UserIdentity struct {
	UserID int64
	Role   UserRole
}

type UserProfile struct {
	UsrProfile  Profile
	UsrPurposes []UserPurpose
	UsrSocials  []UserSocialLink
}
