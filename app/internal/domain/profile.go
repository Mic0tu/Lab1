package domain

import "github.com/google/uuid"

type Profile struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Lname string    `json:"lname"`
}
