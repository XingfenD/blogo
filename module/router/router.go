package router

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
)

var loadedConfig config.Config
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

	loadStatic()
	loadHomepage()
	loadCategories()
	loadTags()
	loadArchives()
	loadAbout()
	loadPosts()
	loadNotFound()
	loadAdmin()
	loadRoot()

}

func loadStatic() {
	http.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("Request for /%s from %s", r.URL.Path, r.RemoteAddr)
		fs := http.Dir(loadedConfig.Basic.RootPath + "/static")
		path := r.URL.Path
		loader.Logger.Infof("Opening file /%s", path)
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
}

func loadRoot() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loader.Logger.Infof("RootRouter: Request for %s from %s", r.URL.Path, r.RemoteAddr)

		if r.URL.Path == "/" {
			http.Redirect(w, r, "/homepage/", http.StatusFound)
		} else {
			http.Redirect(w, r, "/404", http.StatusFound)
		}
	})
}
