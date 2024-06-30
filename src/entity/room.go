package entity

import "time"

type Room struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	AllowUsers []string // TODO:#feature

	UserId string // TODO:#feature, add how user create this room

	IsOpen bool // TODO: #feature

	CreateAt time.Time `json:"create_at"`
	CloseAt  time.Time // TODO: #feature
}
