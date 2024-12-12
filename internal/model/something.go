package model

import "time"

type Something struct {
	Id          uint      `json:"id"`
	Description string    `json:"description"`
	UserId      uint      `json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
