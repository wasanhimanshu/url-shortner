package main

import (
	"time"

	"gorm.io/gorm"
)

type CreateUserPayload struct {
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=16"`
}

type CreateShortUrlPayload struct {
	UrlAddr string `json:"url" validate:"required"`
}

type User struct {
	gorm.Model
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type Url struct {
	gorm.Model
	Shorturl    string `json:"shorturl"`
	RedirectUrl string `json:"redirecturl"`
	Clicks      uint   `json:"clicks"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=16"`
}
