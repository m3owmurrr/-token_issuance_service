package main

import (
	"github.com/m3owmurrr/token_server/internal/app"
	"github.com/m3owmurrr/token_server/pkg/config"
)

func main() {
	config.LoadConfig()
	app.NewServer()
}
