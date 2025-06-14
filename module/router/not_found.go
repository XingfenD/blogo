package router

import (
	"net/http"

	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
	"github.com/XingfenD/blogo/module/tpl"
)

func loadNotFound() {
	http.HandleFunc("/404/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		loader.Logger.Infof("Request for /404/ from %s", r.RemoteAddr)

		if tpl.NotFoundTpl == nil {
			loader.Logger.Error("NotFound template is not initialized successfully")
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}
		err = tpl.NotFoundTpl.Execute(w, struct {
			Config config.Config
			Icons  map[string]string
			Path   string
		}{
			Config: loadedConfig,
			Icons:  iconMap,
			Path:   r.URL.Path,
		})
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			loader.Logger.Error(err)
		}
	})
}
