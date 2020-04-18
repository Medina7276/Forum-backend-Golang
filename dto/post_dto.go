package dto

import (
	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

type PostDto struct {
	ID       uuid.UUID `json:"id"`
	ParentID uuid.UUID `json:"parentid"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`

	// int64 Unix ms format
	CreationDate int64           `json:"creationdate"`
	SubofrumID   uuid.UUID       `json:"subforumid"`
	UserID       uuid.UUID       `json:"userid"`
	Comments     []model.Comment `json:"comments"`
}
