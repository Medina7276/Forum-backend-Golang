package model

import uuid "github.com/satori/go.uuid"

type Subforum struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    uuid.UUID `json:"parentid"`
}
