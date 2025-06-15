package router

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/XingfenD/blogo/module/loader"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
	})
}

// 新增认证中间件
// 修改认证中间件
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

	// 原有 JWT 验证逻辑
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

// 修改登录处理函数
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

// 新增密码验证函数（需要实现密码哈希方法）
func checkPasswordHash(password, hash string) bool {
	// 这里需要实现具体的哈希验证逻辑，例如使用 bcrypt
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
	// return true
}
