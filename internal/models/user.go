package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Role string

const (
	RoleGuest Role = "GUEST"
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName    string             `json:"first_name" bson:"first_name"`
	LastName     string             `json:"last_name" bson:"last_name"`
	Email        string             `json:"email" bson:"email"`
	PhoneNumber  string             `json:"phone_number" bson:"phone_number"`
	PasswordHash string             `json:"-" bson:"password_hash"`
	Role         Role               `json:"role" bson:"role"`
	Balance      float64            `json:"balance" bson:"balance"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

type UserRegistration struct {
	FirstName   string `json:"first_name" binding:"required,min=2"`
	LastName    string `json:"last_name" binding:"required,min=2"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
