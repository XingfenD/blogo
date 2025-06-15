package router

import (
	"fmt"
	"net/http"

	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
	sqlite_db "github.com/XingfenD/blogo/module/sqlite"
	"github.com/XingfenD/blogo/module/tpl"
)

func loadArchives() {
	http.HandleFunc("/archives/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for /archives from %s", r.RemoteAddr)

		err := tpl.ArchiveTpl.Execute(w, struct {
			Config      config.Config
			Icons       map[string]string
			ArchiveMeta ArchiveMeta
		}{
			Config: loadedConfig,
			Icons:  iconMap,
			ArchiveMeta: ArchiveMeta{
				Categories: sqlite_db.GetCategoryList(false),
				Tags:       sqlite_db.GetTagList(false),
				ArticlesOrderByYear: func() map[string][]Terms {
					articles := sqlite_db.GetArticleList()
					yearMap := make(map[string][]Terms)
					for _, a := range articles {
						yearMap[a.Year] = append(yearMap[a.Year], Terms{
							Name: a.Title,
							Url:  fmt.Sprintf("posts/%s", a.DirName),
							Time: a.CreateDate,
						})
					}
					return yearMap
				}(),
			},
		})
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			loader.Logger.Error(err)
		}
	})
}
