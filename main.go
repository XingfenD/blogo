package main

import (
	"github.com/XingfenD/blogo/modules/config"
	"github.com/XingfenD/blogo/modules/router"
)

func main() {
	loaded_config := config.LoadConfig()
	router.StartServer(loaded_config)
}
