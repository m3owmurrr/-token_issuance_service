package services

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/m3owmurrr/token_server/internal/models"
	"github.com/m3owmurrr/token_server/internal/repository"
	"github.com/m3owmurrr/token_server/pkg/config"
	"github.com/m3owmurrr/token_server/pkg/utils"
)

type TokenService struct {
	store repository.Repository
}

func NewTokenService(store repository.Repository) *TokenService {
	return &TokenService{
		store: store,
	}
}

func (ts *TokenService) CreateTokens(user models.User, ip string) (*models.TokensPair, error) {
	tokens, err := utils.GenerateJWTTokens(user.GUID, ip)
	if err != nil {
		return nil, err
	}

	if err := ts.store.PutToken(tokens.GUID, tokens.Refresh); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (ts *TokenService) ValidateTokens(tokens models.TokensRequest) error {

	aTknClaims := jwt.MapClaims{}
	atkn, err := jwt.ParseWithClaims(tokens.Access, aTknClaims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})
	if err != nil {
		return err
	}

	rtknDecode, err := utils.DecodeBase64(tokens.Refresh)
	if err != nil {
		return err
	}

	rTknClaims := jwt.MapClaims{}
	rtkn, err := jwt.ParseWithClaims(rtknDecode, rTknClaims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})
	if err != nil {
		return err
	}

	// Если токены не валидны, то изменялось их содержимое
	if !(atkn.Valid && rtkn.Valid) {
		return err
	}

	// Eсли разные id, то токены не парные
	if aTknClaims["tknid"] != rTknClaims["tknid"] {
		return err
	}

	uuidNotString, _ := uuid.Parse(rTknClaims["tknid"].(string))
	// Если токена нет в бд, то токены не наши
	if err := ts.store.GetToken(uuidNotString); err != nil {
		return err
	}

	return nil
}

func (ts *TokenService) ValidateIP(rtoken, ip string) {
	// ничего непроверям, т.к. перед проверкой ip обязательно должна быть проверка токенов, иначе нет смысла
	tknClaims := jwt.MapClaims{}
	rtknDecode, _ := utils.DecodeBase64(rtoken)
	jwt.ParseWithClaims(rtknDecode, tknClaims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})

	if tknClaims["ip"] != ip {
		utils.SendWarningMail()
	}
}

func (ts *TokenService) CreateTokensRefresh(token, ip string) (*models.TokensPair, error) {
	tknClaims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, tknClaims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})

	guidUser, err := uuid.Parse(tknClaims["sub"].(string))
	if err != nil {
		return nil, err
	}

	tmpTkns, err := ts.CreateTokens(models.User{GUID: guidUser}, ip)
	if err != nil {
		return nil, err
	}

	guidTkn, err := uuid.Parse(tknClaims["tknid"].(string))
	if err != nil {
		return nil, err
	}

	if err := ts.store.DeleteToken(guidTkn); err != nil {
		return nil, err
	}

	return tmpTkns, nil
}
