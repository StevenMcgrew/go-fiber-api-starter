package models

import "time"

type Notification struct {
	Id          uint      `json:"id" form:"id"`
	TextContent string    `json:"textContent" form:"textContent"`
	HasViewed   bool      `json:"hasViewed" form:"hasViewed"`
	UserId      uint      `json:"userId" form:"userId"`
	CreatedAt   time.Time `json:"createdAt" form:"createdAt"`
}
