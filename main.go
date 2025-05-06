package main

import (
	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
	"github.com/XingfenD/blogo/module/router"
)

func main() {
	loaded_config := config.LoadConfig()
	loader.LoadLogger(loaded_config.Basic.LogLevel)
	router.StartServer(loaded_config)
}
