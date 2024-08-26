package models

import "github.com/google/uuid"

type TokensPair struct {
	GUID    uuid.UUID
	Access  string
	Refresh string
}

type Tokens struct {
	Access  string
	Refresh string
}

type TokensResponse struct {
	Access string `json:"access_token"`
}

type TokensRequest struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

type TokensDB struct {
	GUID    uuid.UUID
	Refresh string
}
