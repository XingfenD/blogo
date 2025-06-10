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
)

// StartServer 初始化并启动 HTTP 服务器
func StartServer(loaded_config config.Config) {
	loadRouter(loaded_config)

	server := &http.Server{
		Addr: "localhost:" + strconv.Itoa(loaded_config.Basic.Port2listen),
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

func loadRouter(loaded_config config.Config) {
	icon_map, _ := loader.LoadIcons(loaded_config.Basic.RootPath + "/static/icon")
	/* the static files */
	http.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for %s from %s", r.URL.Path, r.RemoteAddr)
		fs := http.Dir(loaded_config.Basic.RootPath + "/static")
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		loader.Logger.Infof("Request for / from %s", r.RemoteAddr)
		funcMap := template.FuncMap{
			"date": func(format string) string {
				return time.Now().Format(format)
			},
		}

		t := template.New("index.html").Funcs(funcMap)

		t, err = t.ParseFiles(
			loaded_config.Basic.RootPath+"/template/page/index.html",
			loaded_config.Basic.RootPath+"/template/layout/footer.html",
			loaded_config.Basic.RootPath+"/template/layout/sidebar.html",
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
			Config: loaded_config,
			Icons:  icon_map,
		})
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			loader.Logger.Error(err)
		}
	})
	http.HandleFunc("/category/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for /category/ from %s", r.RemoteAddr)
	})
	http.HandleFunc("/about/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for /about/ from %s", r.RemoteAddr)
	})
}
