package model

import uuid "github.com/satori/go.uuid"

type Comment struct {
	ID           uuid.UUID `json:"id"`
	PostID       uuid.UUID `json:"postid"`
	UserID       uuid.UUID `json:"userid"`
	Content      string    `json:"content"`
	CreationDate int64     `json:"creationdate"`
}
