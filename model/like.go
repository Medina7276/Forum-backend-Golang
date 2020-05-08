package model

import uuid "github.com/satori/go.uuid"

type Like struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"userid"`
	PostID   uuid.UUID `json:"postid"`
	IsUpVote bool      `json:"isupvote"`
}
