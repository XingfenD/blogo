package main

import (
	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/router"
)

func main() {
	loaded_config := config.LoadConfig()
	router.StartServer(loaded_config)
}
