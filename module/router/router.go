package router

import (
	"context"
	"log"
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
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func loadRouter(loaded_config config.Config) {
	icon_map, _ := loader.LoadIcons()
	loader.LoadPages()
	indexHTML := func(w http.ResponseWriter, r *http.Request) {
		funcMap := template.FuncMap{
			"date": func(format string) string {
				return time.Now().Format(format)
			},
		}

		t := template.New("index.html").Funcs(funcMap)
		path, err := os.Getwd()
		if err != nil {
			http.Error(w, "Failed to get working directory", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		t, err = t.ParseFiles(
			path+"/template/index.html",
			path+"/template/layout/footer.html",
			path+"/template/layout/sidebar.html",
		)
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = t.Execute(w, struct {
			Config config.Config
			Icons  map[string]string
		}{
			Config: loaded_config,
			Icons:  icon_map,
		})
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			log.Println(err)
		}
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHTML)
	http.HandleFunc("/index.html", indexHTML)
}
