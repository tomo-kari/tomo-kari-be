package entity

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userRole string

const (
	ADMIN   userRole = "admin"
	PARTNER userRole = "partner"
	USER    userRole = "user"
)

type User struct {
	ID          uint64    `json:"id"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Password    string    `json:"password"`
	Role        userRole  `json:"role"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Description string    `json:"description"`
	FacebookId  string    `json:"facebookId"`
	GoogleId    string    `json:"googleId"`
	TimeStamp
}

func (*User) Table() string {
	return "Users"
}

func (u *User) HashPassword() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(bytes)
}

func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type BasicInfo struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserTokenData struct {
	ID    uint64 `json:"id"`
	Role  string `json:"role"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type AuthUserResponse struct {
	ID uint64 `json:"id"`
	BasicInfo
	DateOfBirth time.Time `json:"dateOfBirth"`
	Token
}

type CreateUserRequestBody struct {
	BasicInfo
	DateOfBirth      string `json:"dateOfBirth"`
	Password         string `json:"password"`
	TermsOfServiceId int64  `json:"termsOfServiceId"`
}

type LoginUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
