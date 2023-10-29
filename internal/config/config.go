package config

import (
	"log"
	"os"
	"race-weekend-bot/internal/boto"
	"race-weekend-bot/internal/storage/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string          `yaml:"env" env-default:"local"`
	DatabaseEngine string          `yaml:"database-engine" env-default:"postgres"`
	Bot            boto.BotConfig  `yaml:"bot"`
	Postgress      postgres.Config `yaml:"postgres"`
	Api            APIs            `yaml:"api"`
}

type APIs struct {
	F1     string `yaml:"f1" env-required:"true"`
	Motogp string `yaml:"motogp" env-required:"true"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		log.Fatal("config path is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
