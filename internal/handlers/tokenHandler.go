package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/m3owmurrr/token_server/internal/models"
	"github.com/m3owmurrr/token_server/internal/services"
	"github.com/m3owmurrr/token_server/pkg/utils"
)

type TokenHandler struct {
	serv services.Service
}

func NewTokenHandler(serv services.Service) *TokenHandler {
	return &TokenHandler{
		serv: serv,
	}
}

func (th *TokenHandler) GetTokens(w http.ResponseWriter, r *http.Request) {
	var user models.User

	params := r.URL.Query()
	if params.Get("GUID") == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("empty request params")
		return
	}
	user.GUID, _ = uuid.Parse(params.Get("GUID"))

	tkns, err := th.serv.CreateTokens(user, strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cannot create tokens: %s", err)
		return
	}

	body, err := json.Marshal(tkns.Access)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cannot encode token: %s", err)
		return
	}

	cookie := http.Cookie{
		Name:     "RefreshToken",
		Value:    utils.EncodeBase64(tkns.Refresh),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (th *TokenHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var tknRequest models.TokensRequest

	if err := json.NewDecoder(r.Body).Decode(&tknRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("cannot unmarshal request body: %s", err)
		return
	}

	if err := th.serv.ValidateTokens(tknRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("not valid tokens: %s", err)
		fmt.Fprintln(w, "Not valid tokens!")
		return
	}

	th.serv.ValidateIP(tknRequest.Refresh, strings.Split(r.RemoteAddr, ":")[0])

	tkns, err := th.serv.CreateTokensRefresh(tknRequest.Access, r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cannot create tokens: %s", err)
		return
	}

	body, err := json.Marshal(tkns.Access)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cannot encode token: %s", err)
		return
	}

	cookie := http.Cookie{
		Name:     "RefreshToken",
		Value:    utils.EncodeBase64(tkns.Refresh),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
