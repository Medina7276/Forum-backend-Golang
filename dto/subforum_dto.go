package dto

import (
	"git.01.alem.school/qjawko/forum/model"
	uuid "github.com/satori/go.uuid"
)

type SubforumDto struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	ParentID    uuid.UUID        `json:"parentid"`
	Children    []model.Subforum `json:"children"`
	Posts       []model.Post     `json:"posts"`
}
