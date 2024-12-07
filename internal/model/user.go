package model

import "gorm.io/gorm"

// How to declare gorm models: https://gorm.io/docs/models.html

// User struct
type User struct {
	gorm.Model
	Email    string  `gorm:"uniqueIndex;not null" json:"email"`
	Username string  `gorm:"uniqueIndex;not null" json:"username"`
	Password string  `gorm:"not null" json:"password"`
	Type     string  `gorm:"not null;default:user" json:"type"`
	Status   string  `gorm:"not null;default:active" json:"status"`
	Picture  *string `json:"picture"`
}
