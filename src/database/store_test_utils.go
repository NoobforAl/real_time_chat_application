package database

import (
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/entity"
)

func timesEqual(t1, t2 time.Time) bool {
	return t1.Equal(t2)
}

func roomsEqual(r1, r2 entity.Room) bool {
	if r1.Id != r2.Id || r1.Name != r2.Name || r1.Description != r2.Description || r1.UserId != r2.UserId || r1.IsOpen != r2.IsOpen {
		return false
	}
	if !timesEqual(r1.CreateAt, r2.CreateAt) || !timesEqual(r1.CloseAt, r2.CloseAt) {
		return false
	}
	return true
}

func messagesEqual(m1, m2 entity.Message) bool {
	if m1.Id != m2.Id || m1.Content != m2.Content || m1.SenderId != m2.SenderId || m1.RoomId != m2.RoomId {
		return false
	}
	if !timesEqual(m1.Timestamp, m2.Timestamp) {
		return false
	}
	return true
}

func notificationsEqual(n1, n2 entity.Notification) bool {
	if n1.Id != n2.Id || n1.SenderId != n2.SenderId || n1.RoomId != n2.RoomId || n1.Content != n2.Content || n1.ReadNotification != n2.ReadNotification {
		return false
	}
	if !timesEqual(n1.CreateAt, n2.CreateAt) {
		return false
	}
	return true
}

func notificationRoomsEqual(nr1, nr2 entity.NotificationRoom) bool {
	if nr1.Id != nr2.Id || nr1.UserId != nr2.UserId || nr1.RoomId != nr2.RoomId || nr1.Enable != nr2.Enable {
		return false
	}
	if !timesEqual(nr1.CreateAt, nr2.CreateAt) {
		return false
	}
	return true
}
