package config

import (
	"errors"
	"log"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	AccessTokenLifetime  time.Duration `yaml:"access_token_lifetime"`
	RefreshTokenLifetime time.Duration `yaml:"refresh_token_lifetime"`
	S                    ServerConfig  `yaml:"server"`
	DB                   DBConfig      `yaml:"database"`
}

type DBConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
}

type ServerConfig struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

var (
	Cfg    = &Config{}
	JwtKey = []byte(os.Getenv("JWT_KEY"))
)

func LoadConfig() {
	configPath := "./config/config.yaml"

	f, err := os.Open(configPath)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalf("config file does not exist: %s", configPath)
	} else if err != nil {
		log.Fatalf("config file error: %s", err.Error())
	}

	if err = yaml.NewDecoder(f).Decode(Cfg); err != nil {
		log.Fatalf("cannot read config: %s", err.Error())
	}
}
