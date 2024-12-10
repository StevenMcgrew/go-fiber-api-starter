package model

import "gorm.io/gorm"

// How to declare gorm models: https://gorm.io/docs/models.html
// How to declare validator tags: https://pkg.go.dev/github.com/go-playground/validator#hdr-Using_Validator_Tags

// User struct
type User struct {
	gorm.Model
	Email     string      `gorm:"uniqueIndex;not null" validate:"required,contains=@,min=3,max=320" json:"email"`
	Username  string      `gorm:"uniqueIndex;not null" validate:"required,alphanum,min=1,max=50" json:"username"`
	Password  string      `gorm:"not null" validate:"required,min=8,max=128" json:"password"`
	Type      string      `gorm:"not null;default:user" validate:"oneof=user admin" json:"type"`
	Status    string      `gorm:"not null;default:unverified" validate:"oneof=unverified active suspended deleted" json:"status"`
	ImageURL  *string     `json:"image_url"`
	Something []Something `validate:"nostructlevel"`
}
