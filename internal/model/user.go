package model

import "time"

type User struct {
	Id         uint      `json:"id"`
	Email      string    `json:"email"`
	UserName   string    `json:"userName"`
	Password   string    `json:"password"`
	UserType   string    `json:"userType"`
	UserStatus string    `json:"userStatus"`
	ImageUrl   string    `json:"imageUrl"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
