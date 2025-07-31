package models

import "github.com/google/uuid"

type ChirpPost struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}
