package router

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
	sqlite_db "github.com/XingfenD/blogo/module/sqlite"
	"github.com/XingfenD/blogo/module/tpl"
)

func loadTags() {
	http.HandleFunc("/archives/tags/", func(w http.ResponseWriter, r *http.Request) {
		fullPath := r.URL.Path
		suffix := strings.TrimPrefix(
			path.Clean(fullPath),
			path.Clean("/archives/tags/"),
		)
		if suffix == "" {
			loader.Logger.Infof("Request for /archives/tags/ from %s", r.RemoteAddr)
			err := tpl.SectionTpl.Execute(w, struct {
				Config       config.Config
				Icons        map[string]string
				SectionTitle string
				SectionName  string
				SectionCount int
				Terms        []struct {
					Name string
					Url  string
					Time string
				}
			}{
				Config:       loadedConfig,
				Icons:        iconMap,
				SectionTitle: "SECTION",
				SectionName:  "tags",
				SectionCount: len(sqlite_db.GetTagList()),
				Terms: func() []struct {
					Name string
					Url  string
					Time string
				} {
					var terms []struct {
						Name string
						Url  string
						Time string
					}
					for _, tag := range sqlite_db.GetTagList() {
						terms = append(terms, struct {
							Name string
							Url  string
							Time string
						}{
							Name: tag.Name,
							Url:  fmt.Sprintf("archives/tags/%d", tag.Id),
						})
					}
					return terms
				}(),
			})
			if err != nil {
				http.Error(w, "Failed to execute template", http.StatusInternalServerError)
				loader.Logger.Error(err)
				return
			}
		}
		loader.Logger.Infof("Tag path: %s (From %s)", suffix, r.RemoteAddr)
		parts := strings.Split(suffix, "/")
		if len(parts) > 1 {
			tagID := parts[1]
			loader.Logger.Infof("Requested tag ID: %s", tagID)
			tagIDInt, err := strconv.Atoi(tagID)
			if err != nil {
				http.Error(w, "Invalid tag ID", http.StatusBadRequest)
				loader.Logger.Error("Invalid tag ID:", err)
				return
			}
			sectionName, err := sqlite_db.GetTagById(tagIDInt)
			if err != nil {
				http.Error(w, "Failed to get tag name", http.StatusInternalServerError)
				loader.Logger.Error("Failed to get tag name:", err)
				return
			}
			Articles := sqlite_db.GetArticlesByTag(tagIDInt)
			err = tpl.SectionTpl.Execute(w, struct {
				Config       config.Config
				Icons        map[string]string
				SectionTitle string
				SectionName  string
				SectionCount int
				Terms        []struct {
					Name string
					Url  string
					Time string
				}
			}{
				Config:       loadedConfig,
				Icons:        iconMap,
				SectionTitle: "TAGS",
				SectionName:  sectionName,
				SectionCount: len(Articles),
				Terms: func() []struct {
					Name string
					Url  string
					Time string
				} {
					var terms []struct {
						Name string
						Url  string
						Time string
					}
					for _, article := range Articles {
						terms = append(terms, struct {
							Name string
							Url  string
							Time string
						}{
							Name: article.Title,
							Url:  fmt.Sprintf("posts/%s", article.DirName),
							Time: article.CreateDate,
						})
					}
					return terms
				}(),
			})
			if err != nil {
				http.Error(w, "Failed to execute template", http.StatusInternalServerError)
				loader.Logger.Error(err)
				return
			}
		}

	})
}
