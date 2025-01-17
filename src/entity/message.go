package entity

import "time"

type Message struct {
	Id        string    `json:"id"`
	Content   string    `json:"content" validate:"required,max=64,min=1"`
	SenderId  string    `json:"sender_id"`
	RoomId    string    `json:"room_id"`
	Timestamp time.Time `json:"timestamp"`
}
