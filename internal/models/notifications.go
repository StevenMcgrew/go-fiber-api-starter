package models

import "time"

type Notifications struct {
	Id          uint      `json:"id" form:"id"`
	TextContent string    `json:"textContent" form:"textContent"`
	UserId      uint      `json:"userId" form:"userId"`
	CreatedAt   time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" form:"updatedAt"`
}
