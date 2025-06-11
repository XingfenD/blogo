package main

import (
	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
	"github.com/XingfenD/blogo/module/router"
	sqlite_db "github.com/XingfenD/blogo/module/sqlite"
	"github.com/XingfenD/blogo/module/tpl"
)

func main() {
	loaded_config := config.LoadConfig()
	loader.LoadLogger(loaded_config.Basic.LogLevel)
	sqlite_db.InitDB(loaded_config.Basic.RootPath + "/data/blogo_db.db")
	tpl.LoadTemplate(loaded_config.Basic.RootPath)
	router.StartServer(loaded_config)
}
