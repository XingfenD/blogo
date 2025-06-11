package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
	sqlite_db "github.com/XingfenD/blogo/module/sqlite"
)

var loadedConfig config.Config

var funcMap = template.FuncMap{
	"date": func(format string) string {
		return time.Now().Format(format)
	},
}

var iconMap map[string]string

// StartServer 初始化并启动 HTTP 服务器
func StartServer(loaded_config config.Config) {
	loadedConfig = loaded_config
	loadRouter()

	server := &http.Server{
		Addr: "localhost:" + strconv.Itoa(loadedConfig.Basic.Port2listen),
	}

	go func() {
		loader.Logger.Infof("Starting server on http://%s/", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			loader.Logger.Errorf("Could not start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	loader.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		loader.Logger.Error("Server forced to shutdown:", err)
	}

	loader.Logger.Info("Server exiting")
}

func loadRouter() {
	iconMap, _ = loader.LoadIcons(loadedConfig.Basic.RootPath + "/static/icon")

	/* the static files */
	http.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for %s from %s", r.URL.Path, r.RemoteAddr)
		fs := http.Dir(loadedConfig.Basic.RootPath + "/static")
		path := r.URL.Path
		loader.Logger.Infof("Opening file %s", path)
		file, err := fs.Open(path)
		if err != nil {
			http.NotFound(w, r)
			loader.Logger.Error(err)
			return
		}
		defer file.Close()

		if info, _ := file.Stat(); info.IsDir() {
			http.NotFound(w, r)
			loader.Logger.Error(err)
			return
		}

		http.FileServer(fs).ServeHTTP(w, r)
	})))

	/* the main page */
	http.HandleFunc("/homepage.html", func(w http.ResponseWriter, r *http.Request) {
		var err error
		loader.Logger.Infof("Request for /homepage.html from %s", r.RemoteAddr)

		t := template.New("index.html").Funcs(funcMap)

		t, err = t.ParseFiles(
			loadedConfig.Basic.RootPath+"/template/page/index.html",
			loadedConfig.Basic.RootPath+"/template/layout/footer.html",
			loadedConfig.Basic.RootPath+"/template/layout/sidebar.html",
		)
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			loader.Logger.Error(err)
			return
		}
		loader.Logger.Info("Loading template successfully")
		err = t.Execute(w, struct {
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

	http.HandleFunc("/about.html", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for /about.html from %s", r.RemoteAddr)
		t := template.New("about.html").Funcs(funcMap)
		t, err := t.ParseFiles(
			loadedConfig.Basic.RootPath+"/template/page/about.html",
			loadedConfig.Basic.RootPath+"/template/layout/footer.html",
			loadedConfig.Basic.RootPath+"/template/layout/sidebar.html",
			loadedConfig.Basic.RootPath+"/template/layout/article.html",
		)
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			loader.Logger.Error(err)
			return
		}
		loader.Logger.Info("Loading template successfully")
		aboutMeta, err := sqlite_db.GetAboutMeta()
		// loader.Logger.Info(aboutMeta)
		if err != nil {
			http.Error(w, "Failed to get about meta", http.StatusInternalServerError)
			loader.Logger.Error(err)
			return
		}
		err = t.Execute(w, struct {
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

	loadArchives()
	loadPosts()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for %s from %s", r.URL.Path, r.RemoteAddr)
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/homepage.html", http.StatusFound)
		} else {
			http.NotFound(w, r)
		}
	})
}

func loadArchives() {
	archivesMux := http.NewServeMux()
	archivesMux.HandleFunc("/archives/categories/", func(w http.ResponseWriter, r *http.Request) {
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
			t := template.New("section.html").Funcs(funcMap)
			t, err := t.ParseFiles(
				loadedConfig.Basic.RootPath+"/template/page/section.html",
				loadedConfig.Basic.RootPath+"/template/layout/footer.html",
				loadedConfig.Basic.RootPath+"/template/layout/sidebar.html",
			)
			if err != nil {
				http.Error(w, "Failed to parse template", http.StatusInternalServerError)
				loader.Logger.Error(err)
				return
			}
			loader.Logger.Info("Loading template successfully")
			err = t.Execute(w, struct {
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
				http.Error(w, "Invalid category ID", http.StatusBadRequest)
				loader.Logger.Error("Invalid category ID:", err)
				return
			}
			sectionName, err := sqlite_db.GetCateById(catID)
			if err != nil {
				http.Error(w, "Failed to get category name", http.StatusInternalServerError)
				loader.Logger.Error("Failed to get category name:", err)
				return
			}
			t := template.New("section.html").Funcs(funcMap)
			t, err = t.ParseFiles(
				loadedConfig.Basic.RootPath+"/template/page/section.html",
				loadedConfig.Basic.RootPath+"/template/layout/footer.html",
				loadedConfig.Basic.RootPath+"/template/layout/sidebar.html",
			)
			if err != nil {
				http.Error(w, "Failed to parse template", http.StatusInternalServerError)
				loader.Logger.Error(err)
				return
			}
			loader.Logger.Info("Loading template successfully")
			Articles := sqlite_db.GetArticlesByCategory(catID)

			err = t.Execute(w, struct {
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
							Url:  fmt.Sprintf("posts/%s.html", article.DirName),
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

	archivesMux.HandleFunc("/archives/collections/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for %s from %s", r.URL.Path, r.RemoteAddr)
		w.Write(fmt.Appendf(nil, "Request for %s from %s", r.URL.Path, r.RemoteAddr))
	})

	archivesMux.HandleFunc("/archives/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for /archives from %s", r.RemoteAddr)
		t := template.New("archives.html").Funcs(funcMap)
		t, err := t.ParseFiles(
			loadedConfig.Basic.RootPath+"/template/page/archives.html",
			loadedConfig.Basic.RootPath+"/template/layout/footer.html",
			loadedConfig.Basic.RootPath+"/template/layout/sidebar.html",
		)
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			loader.Logger.Error(err)
			return
		}
		loader.Logger.Info("Loading template successfully")
		err = t.Execute(w, struct {
			Config     config.Config
			Icons      map[string]string
			Categories []struct {
				Name string
				Id   int
				Time string
			}
			Collections []struct {
				Name string
				Id   int
			}
		}{
			Config:      loadedConfig,
			Icons:       iconMap,
			Categories:  sqlite_db.GetCategoryList(),
			Collections: nil,
		})
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			loader.Logger.Error(err)
		}
	})
	http.Handle("/archives/", archivesMux)
}

func loadPosts() {
	postsMux := http.NewServeMux()
	postsMux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		fullPath := r.URL.Path

		suffix := strings.TrimPrefix(
			path.Clean(fullPath),
			path.Clean("/posts/"),
		)
		if suffix == "" {
			loader.Logger.Infof("Request for /posts/ from %s", r.RemoteAddr)
			t := template.New("section.html").Funcs(funcMap)
			t, err := t.ParseFiles(
				loadedConfig.Basic.RootPath+"/template/page/section.html",
				loadedConfig.Basic.RootPath+"/template/layout/footer.html",
				loadedConfig.Basic.RootPath+"/template/layout/sidebar.html",
			)
			if err != nil {
				http.Error(w, "Failed to parse template", http.StatusInternalServerError)
				loader.Logger.Error(err)
				return
			}
			loader.Logger.Info("Loading template successfully")
			posts := sqlite_db.GetArticleList()
			err = t.Execute(w, struct {
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
				SectionTitle: "POSTS",
				SectionName:  "posts",
				SectionCount: len(posts),
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
					for _, article := range posts {
						terms = append(terms, struct {
							Name string
							Url  string
							Time string
						}{
							Name: article.Title,
							Url:  fmt.Sprintf("posts/%s.html", article.DirName),
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
		loader.Logger.Infof("Post path: %s (From %s)", suffix, r.RemoteAddr)
		parts := strings.Split(suffix, "/")
		if len(parts) > 1 {

			// articleDirName := parts[1]

			return
		}

		loader.Logger.Infof("Request for %s from %s", r.URL.Path, r.RemoteAddr)
	})

	http.Handle("/posts/", postsMux)
}
