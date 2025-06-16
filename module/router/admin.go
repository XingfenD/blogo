package router

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/XingfenD/blogo/module/config"
	"github.com/XingfenD/blogo/module/loader"
	sqlite_db "github.com/XingfenD/blogo/module/sqlite"
	"github.com/XingfenD/blogo/module/tpl"
	"github.com/golang-jwt/jwt/v5"
)

// 新增 JWT 配置（需要先在 config 模块添加 JWTSecret 字段）
var jwtSecret = []byte(loadedConfig.Security.JWTSecret)

func loadAdmin() {
	http.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
		fullPath := r.URL.Path
		adminPath := strings.TrimPrefix(fullPath, "/admin/")

		// 添加登录路由处理
		if adminPath == "login" {
			handleAdminLogin(w, r)
			return
		}

		// 其他路由需要 JWT 验证
		if !authMiddleware(w, r) {
			return
		}

		if adminPath == "" {
			handleAdminRender(w, r)
		} else if adminPath != "posts" {
			handleCollePage(adminPath, w, r)
		} else {
			handlePostCollePage(w, r)
		}
	})
}

func authMiddleware(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")

	// 处理 Basic 认证
	if strings.HasPrefix(authHeader, "Basic ") {
		credentials, err := base64.StdEncoding.DecodeString(authHeader[6:])
		if err == nil {
			parts := strings.SplitN(string(credentials), ":", 2)
			if len(parts) == 2 {
				username, password := parts[0], parts[1]
				loader.Logger.Infof("Admin login attempt, username: %s, password: %s", username, password)

				// 验证凭证
				if username == loadedConfig.Security.AdminUser && password == loadedConfig.Security.AdminPass {
					// 生成 JWT
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
						"sub": username,
						"exp": time.Now().Add(time.Hour * 24).Unix(),
					})

					if tokenString, err := token.SignedString(jwtSecret); err == nil {
						// 设置新的 Authorization 头为 Bearer 模式
						w.Header().Set("Authorization", "Bearer "+tokenString)
						return true
					}
				}
			}
		}
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return false
	}

	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		loader.Logger.Warn("Invalid JWT attempt from:", r.RemoteAddr)
		return false
	}
	return true
}

func handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	// 添加 Basic 认证支持
	if authHeader := r.Header.Get("Authorization"); strings.HasPrefix(authHeader, "Basic ") {
		// 通过中间件处理认证
		if authMiddleware(w, r) {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		loader.Logger.Warn("Admin login method not allowed from:", r.RemoteAddr)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// 直接验证配置中的凭证
	if username != "admin" || password != loadedConfig.Security.AdminPass {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		loader.Logger.Warnf("Admin login failed for user: %s", username)
		return
	}

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		loader.Logger.Error("JWT generation failed:", err)
		return
	}

	// 返回 Token（也可以设置为 Cookie）
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"token":"` + tokenString + `"}`))
}

func handleAdminRender(w http.ResponseWriter, r *http.Request) {
	stats := struct {
		TotalPosts      int
		TotalCategories int
		TotalTags       int
	}{
		TotalPosts:      len(sqlite_db.GetArticleList()),
		TotalCategories: len(sqlite_db.GetCategoryList(false)),
		TotalTags:       len(sqlite_db.GetTagList(false)),
	}

	recentPosts := sqlite_db.GetRecentArticles(5)

	err := tpl.AdminTpl.Execute(w, struct {
		Config        config.Config
		Icons         map[string]string
		DashBoardMeta struct {
			Stats       interface{}
			RecentPosts []sqlite_db.ArticleListItem
		}
	}{
		Config: loadedConfig,
		Icons:  iconMap,
		DashBoardMeta: struct {
			Stats       interface{}
			RecentPosts []sqlite_db.ArticleListItem
		}{
			Stats:       stats,
			RecentPosts: recentPosts,
		},
	})
	if err != nil {
		http.Error(w, "Failed to execute admin template", http.StatusInternalServerError)
		loader.Logger.Error("Admin template execution failed:", err)
	}
}

func handleCollePage(collectionName string, w http.ResponseWriter, r *http.Request) {
	err := tpl.ColleTpl.Execute(w, struct {
		Config         config.Config
		Icons          map[string]string
		ColleTableMeta struct {
			Title          string
			CollectionList []sqlite_db.CollectionListItem
		}
	}{
		Config: loadedConfig,
		Icons:  iconMap,
		ColleTableMeta: struct {
			Title          string
			CollectionList []sqlite_db.CollectionListItem
		}{
			Title: collectionName,
			CollectionList: func() []sqlite_db.CollectionListItem {
				if collectionName == "categories" {
					return sqlite_db.GetCategoryList(true)
				} else if collectionName == "tags" {
					return sqlite_db.GetTagList(true)
				} else {
					return nil
				}
			}(),
		},
	})
	if err != nil {
		http.Error(w, "Failed to execute collection template", http.StatusInternalServerError)
		loader.Logger.Error("Collection template execution failed:", err)
	}
}

func handlePostCollePage(w http.ResponseWriter, r *http.Request) {
	err := tpl.PostTableTpl.Execute(w, struct {
		Config           config.Config
		Icons            map[string]string
		ArticleTableMeta []sqlite_db.ArticleListItem
	}{
		Config:           loadedConfig,
		Icons:            iconMap,
		ArticleTableMeta: sqlite_db.GetArticleList(),
	})
	if err != nil {
		http.Error(w, "Failed to execute post collection template", http.StatusInternalServerError)
		loader.Logger.Error("Post collection template execution failed:", err)
	}
}
