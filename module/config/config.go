package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port2listen int `toml:"port2listen"`
	User        struct {
		Name        string `toml:"name"`
		Title       string `toml:"title"`
		Greeting    string `toml:"greeting"`
		Description string `toml:"description"`
		HomeUrl     string `toml:"home_url"`
		AvatarUrl   string `toml:"avatar_url"`
	}
	Sidebar struct {
		Emoji       string `toml:"emoji"`
		SocialLinks []struct {
			Name string `toml:"name"`
			Url  string `toml:"url"`
			Icon string `toml:"icon"`
		} `toml:"social_links"`
	}
	Footer struct {
		CustomText string `toml:"custom_text"`
	}
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
