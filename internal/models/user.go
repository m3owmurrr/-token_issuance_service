package models

import "github.com/google/uuid"

type User struct {
	GUID uuid.UUID `json:"guid"`
	// Some more data
}
