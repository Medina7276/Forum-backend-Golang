package model

import uuid "github.com/satori/go.uuid"

type Post struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`

	// int64 Unix ms format
	CreationDate int64     `json:"creationdate"`
	SubforumID   uuid.UUID `json:"subforumid"`
	UserID       uuid.UUID `json:"userid"`
	// ImageURL     string    `json:"imageurl"`
}
