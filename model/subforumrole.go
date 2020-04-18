package model

import uuid "github.com/satori/go.uuid"

type SubforumRole struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userid"`
	Role       role      `json:"role"`
	SubforumID uuid.UUID `json:"subforumid"`
}

type role int

const (
	ROLE_USER = iota
	ROLE_MODER
	ROLE_ADMIN
)
