package router

import (
	"net/http"

	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
	sqlite_db "github.com/XingfenD/blogo/module/sqlite"
	"github.com/XingfenD/blogo/module/tpl"
	"github.com/russross/blackfriday/v2"
)

func loadAbout() {
	http.HandleFunc("/about/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for /about/ from %s", r.RemoteAddr)

		aboutMeta, err := sqlite_db.GetAboutMeta()
		// loader.Logger.Info(aboutMeta)
		if err != nil {
			http.Error(w, "Failed to get about meta", http.StatusInternalServerError)
			loader.Logger.Error(err)
			return
		}

		/* parse the markdown to html */
		markdown := aboutMeta.Content
		html := blackfriday.Run([]byte(markdown), blackfriday.WithExtensions(
			0|blackfriday.AutoHeadingIDs|
				blackfriday.FencedCode|
				blackfriday.Tables|
				blackfriday.Strikethrough|
				blackfriday.DefinitionLists),
		)
		aboutMeta.Content = string(html)

		err = tpl.PostTpl.Execute(w, struct {
			Config  config.Config
			Icons   map[string]string
			Article sqlite_db.ArticleMeta
		}{
			Config:  loadedConfig,
			Icons:   iconMap,
			Article: *aboutMeta,
		})
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			loader.Logger.Error(err)
		}
	})
}
