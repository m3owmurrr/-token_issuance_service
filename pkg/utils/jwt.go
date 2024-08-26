package utils

import (
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/m3owmurrr/token_server/internal/models"
	"github.com/m3owmurrr/token_server/pkg/config"
)

func GenerateJWTTokens(GUID uuid.UUID, ip string) (*models.TokensPair, error) {
	// payload
	tknGUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	aClaims := jwt.MapClaims{
		"tknid": tknGUID,
		"sub":   GUID,
		"ip":    ip,
		"exp":   time.Now().Add(config.Cfg.AccessTokenLifetime).Unix(),
	}
	aToken := jwt.NewWithClaims(jwt.SigningMethodHS512, aClaims)
	atkn, err := aToken.SignedString([]byte(config.JwtKey))
	if err != nil {
		return nil, err
	}

	rClaims := jwt.MapClaims{
		"tknid": tknGUID,
		"ip":    ip,
		"exp":   time.Now().Add(config.Cfg.RefreshTokenLifetime).Unix(),
	}
	rToken := jwt.NewWithClaims(jwt.SigningMethodHS512, rClaims)
	rtkn, err := rToken.SignedString([]byte(config.JwtKey))
	if err != nil {
		return nil, err
	}

	tmp := &models.TokensPair{
		GUID:    tknGUID,
		Access:  atkn,
		Refresh: rtkn,
	}

	return tmp, nil
}

func EncodeBase64(rtknString string) string {
	return (base64.StdEncoding.EncodeToString([]byte(rtknString)))
}

func DecodeBase64(encodedStr string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedStr)
	return string(decoded), err
}
