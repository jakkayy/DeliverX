package model

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleCustomer Role = "CUSTOMER"
	RoleDriver   Role = "DRIVER"
	RoleAdmin    Role = "ADMIN"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string     `gorm:"size:100;not null" json:"name"`
	Phone        string     `gorm:"size:20;uniqueIndex;not null" json:"phone"`
	Email        *string    `gorm:"size:100;uniqueIndex" json:"email,omitempty"`
	PasswordHash string     `gorm:"column:password_hash;not null" json:"-"`
	Role         Role       `gorm:"size:20;not null" json:"role"`
	ProfileImage *string    `gorm:"column:profile_image;size:500" json:"profile_image,omitempty"`
	FcmToken     *string    `gorm:"column:fcm_token;size:500" json:"-"`
	IsActive     bool       `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (User) TableName() string { return "users" }
