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

func loadCategories() {
	http.HandleFunc("/archives/categories/", func(w http.ResponseWriter, r *http.Request) {
		// 获取完整请求路径
		fullPath := r.URL.Path

		// 截取前缀后的剩余路径
		suffix := strings.TrimPrefix(
			path.Clean(fullPath),
			path.Clean("/archives/categories/"),
		)

		// 处理空路径的情况（当访问 /archives/categories/ 时）
		if suffix == "" {
			loader.Logger.Infof("Request for /archives/categories/ from %s", r.RemoteAddr)
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
				SectionName:  "categories",
				SectionCount: len(sqlite_db.GetCategoryList()),
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
					for _, category := range sqlite_db.GetCategoryList() {
						terms = append(terms, struct {
							Name string
							Url  string
							Time string
						}{
							Name: category.Name,
							Url:  fmt.Sprintf("archives/categories/%d", category.Id),
							Time: category.Time,
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
			return
		}

		loader.Logger.Infof("Category path: %s (From %s)", suffix, r.RemoteAddr)

		parts := strings.Split(suffix, "/")
		if len(parts) > 1 {
			categoryID := parts[1]
			loader.Logger.Infof("Requested category ID: %s", categoryID)
			catID, err := strconv.Atoi(categoryID)
			if err != nil {
				http.Redirect(w, r, "/404", http.StatusFound)
				loader.Logger.Error("Invalid category ID:", err)
				return
			}
			sectionName, err := sqlite_db.GetCateById(catID)
			if err != nil {
				http.Redirect(w, r, "/404", http.StatusFound)
				loader.Logger.Error("Failed to get category name:", err)
				return
			}
			Articles := sqlite_db.GetArticlesByCategory(catID)

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
				SectionTitle: "CATEGORIES",
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
			return
		}

		http.NotFound(w, r)
	})
}
