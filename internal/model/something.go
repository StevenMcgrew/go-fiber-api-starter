package model

import "gorm.io/gorm"

// How to declare gorm models: https://gorm.io/docs/models.html
// How to declare validator tags: https://pkg.go.dev/github.com/go-playground/validator#hdr-Using_Validator_Tags

// Product struct
type Something struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	UserID      uint
}
