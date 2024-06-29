package entity

import "time"

type Room struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	AllowUsers []string // TODO:#facture

	UserId string // TODO:#facture, add how user create this room

	IsOpen bool // TODO: #facture

	CreateAt time.Time `json:"create_at"`
	CloseAt  time.Time // TODO: #facture
}
