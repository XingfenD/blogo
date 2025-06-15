package router

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
	sqlite_db "github.com/XingfenD/blogo/module/sqlite"
	"github.com/XingfenD/blogo/module/tpl"
	"github.com/russross/blackfriday/v2"
)

func loadPosts() {
	http.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		fullPath := r.URL.Path

		suffix := strings.TrimPrefix(
			path.Clean(fullPath),
			path.Clean("/posts/"),
		)
		if suffix == "" {
			loader.Logger.Infof("Request for /posts/ from %s", r.RemoteAddr)
			postHandler(w, r)
			return
		}

		loader.Logger.Infof("Post path: %s (From %s)", suffix, r.RemoteAddr)
		parts := strings.Split(suffix, "/")
		if len(parts) > 1 {
			postDetailHandler(parts[1], w, r)
			return
		}
	})
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	posts := sqlite_db.GetArticleList()

	err := tpl.SectionTpl.Execute(w, struct {
		Config      config.Config
		Icons       map[string]string
		SectionMeta sectionMeta
	}{
		Config: loadedConfig,
		Icons:  iconMap,
		SectionMeta: sectionMeta{
			SectionTitle: "POSTS",
			SectionName:  "posts",
			SectionCount: len(posts),
			SectionTerms: createPostTerms(posts),
		},
	})
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		loader.Logger.Error(err)
	}
}

func postDetailHandler(dirName string, w http.ResponseWriter, r *http.Request) {
	loader.Logger.Infof("Requested article dir name: %s", dirName)
	article, err := sqlite_db.GetArticleMetaByDir(dirName)

	if err != nil {
		http.Redirect(w, r, "/404", http.StatusFound)
		loader.Logger.Error("Failed to get article:", err)
		return
	}

	/* parse the markdown to html */
	markdown := article.Content
	html := blackfriday.Run([]byte(markdown), blackfriday.WithExtensions(
		0|blackfriday.AutoHeadingIDs|
			blackfriday.FencedCode|
			blackfriday.Tables|
			blackfriday.Strikethrough|
			blackfriday.DefinitionLists),
	)
	article.Content = string(html)

	err = tpl.PostTpl.Execute(w, struct {
		Config  config.Config
		Icons   map[string]string
		Article sqlite_db.ArticleMeta
	}{
		Config:  loadedConfig,
		Icons:   iconMap,
		Article: *article,
	})
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		loader.Logger.Error(err)
	}
}

func createPostTerms(posts []sqlite_db.ArticleListItem) []Terms {
	var terms []Terms
	for _, article := range posts {
		terms = append(terms, Terms{
			Name: article.Title,
			Url:  fmt.Sprintf("posts/%s.html", article.DirName),
			Time: article.CreateDate,
		})
	}
	return terms
}
