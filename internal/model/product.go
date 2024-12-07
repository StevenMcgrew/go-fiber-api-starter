package model

import "gorm.io/gorm"

// How to declare gorm models: https://gorm.io/docs/models.html

// Product struct
type Product struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	Amount      int    `gorm:"not null" json:"amount"`
}
