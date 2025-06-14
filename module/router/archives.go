package router

import (
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
			Config     config.Config
			Icons      map[string]string
			Categories []struct {
				Name string
				Id   int
				Time string
			}
			Tags []struct {
				Name string
				Id   int
			}
			ArticlesOrderByYear map[string][]struct {
				Title   string
				DirName string
				Time    string
			}
		}{
			Config:     loadedConfig,
			Icons:      iconMap,
			Categories: sqlite_db.GetCategoryList(),
			Tags:       sqlite_db.GetTagList(),
			ArticlesOrderByYear: func() map[string][]struct {
				Title   string
				DirName string
				Time    string
			} {
				articles := sqlite_db.GetArticleList()
				yearMap := make(map[string][]struct {
					Title   string
					DirName string
					Time    string
				})
				for _, a := range articles {
					yearMap[a.Year] = append(yearMap[a.Year], struct {
						Title   string
						DirName string
						Time    string
					}{
						Title:   a.Title,
						DirName: a.DirName,
						Time:    a.CreateDate,
					})
				}
				return yearMap
			}(),
		})
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			loader.Logger.Error(err)
		}
	})
}
