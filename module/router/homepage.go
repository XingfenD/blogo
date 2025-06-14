package router

import (
	"net/http"

	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
	"github.com/XingfenD/blogo/module/tpl"
)

func loadHomepage() {
	http.HandleFunc("/homepage/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		loader.Logger.Infof("Request for /homepage/ from %s", r.RemoteAddr)

		err = tpl.IndexTpl.Execute(w, struct {
			Config config.Config
			Icons  map[string]string
		}{
			Config: loadedConfig,
			Icons:  iconMap,
		})
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			loader.Logger.Error(err)
		}
	})
}
