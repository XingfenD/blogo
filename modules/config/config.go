package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port2listen int `toml:"port2listen"`
}

func LoadConfig() Config {
	var config Config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Printf("Config loaded: %+v", config)

	if !checkConfig() {
		log.Fatalf("Config check failed")
	}

	return config
}

func checkConfig() bool {

	return true
}
