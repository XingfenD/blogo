package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Security struct {
		JWTSecret string `toml:"jwt_secret"`
		AdminUser string `toml:"admin_user"`
		AdminPass string `toml:"admin_pass"`
	}
	Basic struct {
		BaseUrl     string `toml:"base_url"`
		LogLevel    int    `toml:"log_level"`
		Port2listen int    `toml:"port2listen"`
		RootPath    string `toml:"root_path"`
	}
	User struct {
		Name        string `toml:"name"`
		Title       string `toml:"title"`
		Greeting    string `toml:"greeting"`
		Description string `toml:"description"`
		HomeUrl     string `toml:"home_url"`
		AvatarUrl   string `toml:"avatar_url"`
	}
	Sidebar struct {
		Emoji             string `toml:"emoji"`
		EnableDivider     bool   `toml:"enable_divider"`
		EnableSocialLinks bool   `toml:"enable_social_links"`
		EnableAdmin       bool   `toml:"enable_admin"`
		HomeIcon          string `toml:"home_icon"`
		SocialLinks       []struct {
			Name string `toml:"name"`
			Url  string `toml:"url"`
			Icon string `toml:"icon"`
		} `toml:"social_links"`
	}
	Footer struct {
		CustomText string `toml:"custom_text"`
	}
	Archive struct {
		ShowTags     bool `toml:"show_tags"`
		ShowTimeline bool `toml:"show_timeline"`
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
		log.Fatal("Config check failed")
	}

	return config
}

func checkConfig() bool {

	return true
}
