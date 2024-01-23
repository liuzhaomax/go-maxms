package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID        string `gorm:"index:idx_user_id;unique;varchar(50);not null"`
	Username      string `gorm:"index:idx_username;unique;varchar(30);not null"`
	Password      string `gorm:"varchar(30);not null"`
	Mobile        string `gorm:"index:idx_mobile;unique;varchar(14);not null"`
	Email         string `gorm:"index:idx_email;unique;varchar(30);not null"`
	EmailVerified bool   `gorm:"boolean;not null;default:false"`
}
