package model

import "time"

type Something struct {
	Id          uint      `json:"id" form:"id"`
	Description string    `json:"description" form:"description"`
	UserId      uint      `json:"userId" form:"userId"`
	CreatedAt   time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" form:"updatedAt"`
}
