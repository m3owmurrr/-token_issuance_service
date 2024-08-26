package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/m3owmurrr/token_server/internal/handlers"
	"github.com/m3owmurrr/token_server/internal/repository/pgstore"
	"github.com/m3owmurrr/token_server/internal/services"
	"github.com/m3owmurrr/token_server/pkg/config"
	"github.com/m3owmurrr/token_server/pkg/utils"
)

func NewServer() {
	db, err := utils.NewDBConnection(10, 10, 30)
	if err != nil {
		log.Fatal("cannot connect to database")
	}

	tokenStore := pgstore.NewTokenRepository(db)
	tokenService := services.NewTokenService(tokenStore)
	tokenHandler := handlers.NewTokenHandler(tokenService)
	healthHandler := handlers.NewHealthHandler()

	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods(http.MethodGet)
	sub := router.PathPrefix("/api/v1").Subrouter()
	sub.HandleFunc("/tokens", tokenHandler.GetTokens).Methods(http.MethodGet)
	sub.HandleFunc("/tokens/refresh", tokenHandler.RefreshTokens).Methods(http.MethodPost)

	loggedRouter := ghandlers.LoggingHandler(os.Stdout, router)

	addr := config.Cfg.S.Host + ":" + config.Cfg.S.Port
	s := http.Server{
		Addr:         addr,
		Handler:      loggedRouter,
		ReadTimeout:  config.Cfg.S.Timeout,
		WriteTimeout: config.Cfg.S.Timeout,
		IdleTimeout:  config.Cfg.S.IdleTimeout,
	}
	fmt.Printf("server running on port: %s\n", config.Cfg.S.Port)
	s.ListenAndServe()
}
