package dto

import (
	_ "github.com/go-playground/validator/v10"
)

type CreateProfileRequest struct {
	Info        ProfileInfoDTO   `json:"info" validate:"required"`
	SocialLinks []SocialLinkDTO  `json:"social_links" validate:"dive"`
	Purposes    []UserPurposeDTO `json:"purposes" validate:"required,min=1,dive"`
}

type ProfileInfoDTO struct {
	Avatar      string  `json:"avatar" validate:"required"`
	Username    string  `json:"username" validate:"required,min=3,max=30"`
	Name        string  `json:"name" validate:"required"`
	Surname     string  `json:"surname" validate:"required"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Email       string  `json:"email" validate:"required,email"`
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty,e164"`
	Sex         *int    `json:"sex" validate:"omitempty,oneof=0 1"`
	Description string  `json:"description" validate:"max=500"`
	Profession  *string `json:"profession,omitempty"`
}

type SocialLinkDTO struct {
	Type string `json:"type" validate:"required"`
	URL  string `json:"url" validate:"required,url"`
}

type UserPurposeDTO struct {
	Purpose string `json:"purpose" validate:"required,min=5"`
}

type UpdateProfileInfoDTO struct {
	
}
