package entity

import "time"

type User struct {
	Id       string `json:"id"`
	Username string `json:"username" validate:"required,max=64,min=8"`
	Email    string `json:"email" validate:"required,email,max=64,min=8"`
	Password string `json:"password" validate:"required,max=64,min=8"`

	Notification bool `json:"notification"` // TODO: #feature

	CreateAt time.Time `json:"create_at"`
}
