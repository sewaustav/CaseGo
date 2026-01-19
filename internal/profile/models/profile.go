package models

type UserSex int 

const (
	Male UserSex = iota
	Female 
)

type Profile struct {
	ID int64
	UserID int64
	Avatar string
	Username string
	Name string
	Surname string
	Patronomyc *string
	PhoneNumber *string
	CaseCount int
	Sex UserSex
}

type UserSocialLink struct {
	ID     int64
	UserID int64
	Type   string 
	URL    string
}