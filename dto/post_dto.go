package dto

import (
	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

type PostDto struct {
	ID           uuid.UUID       `json:"id"`
	ParentID     uuid.UUID       `json:"parentid"`
	Title        string          `json:"title"`
	Content      string          `json:"content"`
	CreationDate int64           `json:"creationdate"`
	Subforum     *model.Subforum `json:"subforumid"`
	User         *model.User     `json:"userid"`
	Comments     []model.Comment `json:"comments"`
	Likes        []model.Like    `json:"likes"`
}
