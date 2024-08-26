package repository

import (
	"github.com/google/uuid"
)

type Repository interface {
	GetToken(guid uuid.UUID) error
	PutToken(guid uuid.UUID, rtkn string) error
	DeleteToken(guid uuid.UUID) error
}
