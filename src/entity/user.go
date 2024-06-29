package entity

import "time"

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Notification bool `json:"notification"`

	CreateAt time.Time `json:"create_at"`
}
