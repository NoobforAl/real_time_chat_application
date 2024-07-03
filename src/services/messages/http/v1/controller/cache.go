package controller

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/gofiber/contrib/websocket"
)

type userSession struct {
	conn         *websocket.Conn
	userId       string
	username     string
	notification bool
}

type roomSession struct {
	users      map[string]*userSession // key is user_id
	lastUpdate time.Time
	mux        sync.Mutex
	log        contract.Logger
}

type wsSessionStore struct {
	rooms map[string]*roomSession
	mux   sync.RWMutex
}

var wsSes *wsSessionStore

func init() {
	wsSes = &wsSessionStore{
		rooms: make(map[string]*roomSession),
		mux:   sync.RWMutex{},
	}
}

var (
	prefix = []string{
		"message-section--",      //0
		"notification-section--", //1
	}
)

func (ws *wsSessionStore) getRoom(roomId string, typeRoom int) (*roomSession, error) {
	ws.mux.RLock()
	defer ws.mux.RUnlock()

	key := prefix[typeRoom] + roomId
	val, ok := ws.rooms[key]
	if !ok {
		return nil, fmt.Errorf("not found room!")
	}
	return val, nil
}

func (rs *roomSession) isOn(txt string) bool {
	txt = strings.Trim(strings.ToLower(txt), " ")
	if txt == "on" {
		return true
	} else if txt == "off" {
		return false
	}

	rs.log.Error("in valid error !")
	return true
}

func getRoomSession(roomId string, typeRoom int, log contract.Logger) *roomSession {
	wsSes.mux.Lock()
	defer wsSes.mux.Unlock()

	usersSession := wsSes.rooms

	key := prefix[typeRoom] + roomId
	if roomSes, ok := usersSession[key]; ok {
		return roomSes
	}

	usersSession[key] = &roomSession{
		users:      make(map[string]*userSession),
		lastUpdate: time.Now(),
		mux:        sync.Mutex{},
		log:        log,
	}
	return usersSession[key]
}

func (rs *roomSession) addUser(user *userSession) {
	rs.mux.Lock()
	defer rs.mux.Unlock()

	rs.users[user.userId] = user
}

func (rs *roomSession) delUser(userId string) {
	rs.mux.Lock()
	defer rs.mux.Unlock()

	delete(rs.users, userId)
}

func (rs *roomSession) notify(message any, senderId string) {
	rs.mux.Lock()
	defer rs.mux.Unlock()

	for userId, user := range rs.users {
		if userId != senderId && user.notification {
			rs.log.Debugf("Sending message from user id (%s) to user id: %s message content: %v", senderId, userId, message)
			err := user.conn.WriteJSON(message)
			if err != nil {
				rs.log.Errorf("Failed to send message to user id: %s, error: %v", userId, err)
				if websocket.IsCloseError(err) {
					rs.log.Warnf("WebSocket connection closed for user id: %s", userId)
					rs.delUser(userId)
					return
				}
			}
		}
	}
}

func (rs *roomSession) turnOffOrOnNotification(ctx context.Context, userId, roomId string) {
	rs.mux.Lock()
	userConn := rs.users[userId]
	rs.mux.Unlock()

	defer func() {
		if err := recover(); err != nil {
			rs.log.Error(err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return

		default:
			_, txt, err := userConn.conn.ReadMessage()
			if websocket.IsCloseError(err) {
				return
			} else if err != nil {
				rs.log.Error(err)
				continue
			}

			userConn.notification = rs.isOn(string(txt))
			userMessageConn, err := wsSes.getRoom(roomId, 0)
			if err == nil {
				userMessageConn.mux.Lock()

				user, ok := userMessageConn.users[userId]
				if ok {
					rs.log.Debugf("turn on/off user notification id: %s  notif: %v", user.userId, userConn.notification)
					user.notification = userConn.notification
				}

				userMessageConn.mux.Unlock()
			}
		}
	}
}
