package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/vizurth/chat-server/internal/postgres"
)

type Config struct {
	Postgres postgres.Config `yaml:"postgres"`
	Port     string          `yaml:"port" env-default:"50052"`
}

func NewConfig() (Config, error) {
	var c Config
	if err := cleanenv.ReadConfig("./config/config.yaml", &c); err != nil {
		fmt.Println(err)
		if err = cleanenv.ReadEnv(&c); err != nil {
			return Config{}, err
		}
	}
	return c, nil
}
