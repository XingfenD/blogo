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
			tagHandler(w)
			return
		}

		loader.Logger.Infof("Tag path: %s (From %s)", suffix, r.RemoteAddr)
		parts := strings.Split(suffix, "/")
		if len(parts) > 1 {
			tagDetailHandler(parts[1], w)
			return
		}

	})
}

func tagHandler(w http.ResponseWriter) {
	tags := sqlite_db.GetTagList(false)
	err := tpl.SectionTpl.Execute(w, struct {
		Config      config.Config
		Icons       map[string]string
		SectionMeta sectionMeta
	}{
		Config: loadedConfig,
		Icons:  iconMap,
		SectionMeta: sectionMeta{
			SectionTitle: "TAGS",
			SectionName:  "tags",
			SectionCount: len(tags),
			SectionTerms: func() []Terms {
				var terms []Terms
				for _, tag := range tags {
					terms = append(terms, Terms{
						Name: tag.ColleName,
						Url:  fmt.Sprintf("archives/tags/%d", tag.ColleId),
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

func tagDetailHandler(tagID string, w http.ResponseWriter) {
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
		Config      config.Config
		Icons       map[string]string
		SectionMeta sectionMeta
		Terms       []struct {
			Name string
			Url  string
			Time string
		}
	}{
		Config: loadedConfig,
		Icons:  iconMap,
		SectionMeta: sectionMeta{
			SectionTitle: "TAGS",
			SectionName:  sectionName,
			SectionCount: len(Articles),
			SectionTerms: createTagTerms(Articles),
		},
	})
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		loader.Logger.Error(err)
		return
	}
}

func createTagTerms(Articles []sqlite_db.ArticleListItem) []Terms {
	var terms []Terms
	for _, article := range Articles {
		terms = append(terms, Terms{
			Name: article.Title,
			Url:  fmt.Sprintf("posts/%s", article.DirName),
			Time: article.CreateDate,
		})
	}
	return terms
}
