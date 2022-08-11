package entity

import (
	"time"
)

type userRole string

const (
	ADMIN   userRole = "admin"
	PARTNER userRole = "partner"
	USER    userRole = "user"
)

type User struct {
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Password    string    `json:"password"`
	Role        userRole  `json:"role"`
	Balance     int       `json:"balance"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Description string    `json:"description"`
	FacebookId  string    `json:"facebookId"`
	GoogleId    string    `json:"googleId"`
}

type BasicInfo struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenData struct {
	ID    uint   `json:"id"`
	Role  string `json:"role"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type AuthUserResponse struct {
	BasicInfo
	DateOfBirth time.Time `json:"dateOfBirth"`
	Token
}

type CreateUserRequestBody struct {
	BasicInfo
	DateOfBirth      string `json:"dateOfBirth"`
	Password         string `json:"password"`
	TermsOfServiceId int    `json:"termsOfServiceId"`
}

type LoginUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
