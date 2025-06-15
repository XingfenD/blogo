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
			cateHandler(w, r)
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
			cateDetailHandler(catID, w, r)
			return
		}

		http.NotFound(w, r)
	})
}

func cateHandler(w http.ResponseWriter, r *http.Request) {
	loader.Logger.Infof("Request for /archives/categories/ from %s", r.RemoteAddr)
	categories := sqlite_db.GetCategoryList(false)

	err := tpl.SectionTpl.Execute(w, struct {
		Config      config.Config
		Icons       map[string]string
		SectionMeta sectionMeta
	}{
		Config: loadedConfig,
		Icons:  iconMap,
		SectionMeta: sectionMeta{
			SectionTitle: "SECTION",
			SectionName:  "categories",
			SectionCount: len(categories),
			SectionTerms: func() []Terms {
				var terms []Terms
				for _, category := range categories {
					terms = append(terms, Terms{
						Name: category.ColleName,
						Url:  fmt.Sprintf("archives/categories/%d", category.ColleId),
					})
				}
				return terms
			}(),
		},
	})
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		loader.Logger.Error(err)
	}
}

func cateDetailHandler(cateId int, w http.ResponseWriter, r *http.Request) {
	sectionName, err := sqlite_db.GetCateById(cateId)
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusFound)
		loader.Logger.Error("Failed to get category name:", err)
		return
	}
	Articles := sqlite_db.GetArticlesByCategory(cateId)

	err = tpl.SectionTpl.Execute(w, struct {
		Config      config.Config
		Icons       map[string]string
		SectionMeta sectionMeta
	}{
		Config: loadedConfig,
		Icons:  iconMap,
		SectionMeta: sectionMeta{
			SectionTitle: "CATEGORIES",
			SectionName:  sectionName,
			SectionCount: len(Articles),
			SectionTerms: func() []Terms {
				var terms []Terms
				for _, article := range Articles {
					terms = append(terms, Terms{
						Name: article.Title,
						Url:  fmt.Sprintf("posts/%s", article.DirName),
						Time: article.CreateDate,
					})
				}
				return terms
			}(),
		},
	})
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		loader.Logger.Error(err)
	}
}
