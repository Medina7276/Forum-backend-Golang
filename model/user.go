package model

import (
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"password,omitempty"`
	AvatarURL string    `json:"avatarurl"`
}
