package models

import "time"

type User struct {
	Id         uint      `json:"id" form:"id"`
	Email      string    `json:"email" form:"email"`
	UserName   string    `json:"userName" form:"userName"`
	Password   string    `json:"password" form:"password"`
	UserType   string    `json:"userType" form:"userType"`
	UserStatus string    `json:"userStatus" form:"userStatus"`
	ImageUrl   string    `json:"imageUrl" form:"imageUrl"`
	CreatedAt  time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt" form:"updatedAt"`
}
