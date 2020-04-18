package model

import uuid "github.com/satori/go.uuid"

type Post struct {
	ID       uuid.UUID `json:"id"`
	ParentID uuid.UUID `json:"parentid"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`

	// int64 Unix ms format
	CreationDate int64     `json:"creationdate"`
	SubofrumID   uuid.UUID `json:"subforumid"`
	UserID       uuid.UUID `json:"userid"`
	// ImageURL     string    `json:"imageurl"`
}
