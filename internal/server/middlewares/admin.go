package middlewares

import (
	"auther/configs"
	"net/http"
)

func AdminTokenMiddleware(cfg *configs.AdminConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("Authorization")
			if !isAdminToken(authToken, cfg.Tokens) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isAdminToken(token string, adminTokens []string) bool {
	for _, adminToken := range adminTokens {
		if token == adminToken {
			return true
		}
	}
	return false
}
