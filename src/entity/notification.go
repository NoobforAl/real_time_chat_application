package entity

import "time"

type Notification struct {
	Id                 string    `json:"id"`
	SenderId           string    `json:"sender_id"`
	RoomId             string    `json:"room_id"`
	NotificationRoomId string    `json:"notification_room_id"`
	Content            string    // TODO: #feature
	CreateAt           time.Time `json:"create_at"`
	ReadNotification   bool      `json:"read_notification"`
}

type NotificationRoom struct {
	Id       string    `json:"id"`
	UserId   string    `json:"user_id"`
	RoomId   string    `json:"room_id"`
	CreateAt time.Time `json:"create_at"`
	Enable   bool      `json:"enable"`
}
