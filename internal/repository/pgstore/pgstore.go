package pgstore

import (
	"crypto/sha256"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/m3owmurrr/token_server/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotExists = errors.New("token with this GUID doesn't exist")

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) GetToken(guid uuid.UUID) error {
	row := r.db.QueryRow("SELECT * FROM tokens WHERE id = $1", guid)

	var tkn models.TokensDB

	if err := row.Scan(&tkn.GUID, &tkn.Refresh); err == sql.ErrNoRows {
		return err
	}

	return nil
}

func (tr *TokenRepository) PutToken(guid uuid.UUID, rtkn string) error {

	// Необходимо, т.к. в исходном виде refresh-токен превышает 72 байта (максимальное значение для bcrypt)
	hash := sha256.New()
	hash.Write([]byte(rtkn))
	hashedRtkn := hash.Sum(nil)

	hashedToken, err := bcrypt.GenerateFromPassword(hashedRtkn, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = tr.db.Exec("INSERT INTO tokens (GUID, refresh_token) VALUES ($1,$2)", guid, string(hashedToken))
	if err != nil {
		return err
	}

	return nil
}

func (tr *TokenRepository) DeleteToken(guid uuid.UUID) error {
	_, err := tr.db.Exec("DELETE FROM tokens WHERE GUID = $1", guid)
	if err != nil {
		return err
	}

	return nil
}
