package router

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/template"
	"time"

	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
	sqlite_db "github.com/XingfenD/blogo/module/sqlite"
)

type page_data struct {
	Config  config.Config
	Icons   map[string]string
	Article sqlite_db.ArticleMeta
}

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
		err = t.Execute(w, page_data{
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
		err = t.Execute(w, page_data{
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

	loadAPI()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for %s from %s", r.URL.Path, r.RemoteAddr)
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/homepage.html", http.StatusFound)
		} else {
			http.NotFound(w, r)
		}
	})
}

func loadPost() {
	postMux := http.NewServeMux()
	postMux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for %s from %s", r.URL.Path, r.RemoteAddr)
	})

	postMux.HandleFunc("/post/{blog_name}.html", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for /post/{blog_name}.html from %s", r.RemoteAddr)
		t := template.New("post.html").Funcs(funcMap)
		t, err := t.ParseFiles(
			loadedConfig.Basic.RootPath+"/template/page/post.html",
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
		err = t.Execute(w, page_data{
			Config: loadedConfig,
			Icons:  iconMap,
			Article: sqlite_db.ArticleMeta{
				Title:        "Sample Post Title",
				CreateDate:   time.Now().Format("2006-01-02"),
				LastModified: time.Now().Format("2006-01-02"),
				Description:  "This is a sample post description.",
				Content:      "This is the content of the sample post.",
			},
		})
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			loader.Logger.Error(err)
		}
	})

	http.Handle("/post/", postMux)
}

func loadArchives() {
	archivesMux := http.NewServeMux()
	archivesMux.HandleFunc("/archives/category/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for %s from %s", r.URL.Path, r.RemoteAddr)
	})

	archivesMux.HandleFunc("/archives/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for /archives/ from %s", r.RemoteAddr)
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
		err = t.Execute(w, page_data{
			Config: loadedConfig,
			Icons:  iconMap,
		})
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			loader.Logger.Error(err)
		}
	})
	http.Handle("/archives/", archivesMux)
}

func loadAPI() {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.Handle("/api/", apiMux)
}
