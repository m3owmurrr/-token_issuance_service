package services

import "github.com/m3owmurrr/token_server/internal/models"

type Service interface {
	CreateTokens(user models.User, ip string) (*models.TokensPair, error)
	CreateTokensRefresh(token, ip string) (*models.TokensPair, error)
	ValidateTokens(tokens models.TokensRequest) error
	ValidateIP(rtoken, ip string)
}
