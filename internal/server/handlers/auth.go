package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"net/http"

	"auther/configs"
	"auther/internal/auth"
)

func LoginHandler(cfg *configs.JWTConfig, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user, err := auth.AuthenticateUser(db, credentials.Login, credentials.Password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, refreshToken, err := auth.GenerateToken(user, cfg)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		if err := auth.SaveRefreshToken(db, user.ID, refreshToken); err != nil {
			http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token, "refresh_token": refreshToken})
	}
}

func RefreshTokenHandler(cfg *configs.JWTConfig, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var refreshTokenStruct struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := json.NewDecoder(r.Body).Decode(&refreshTokenStruct); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		token, err := auth.ParseRefreshToken(refreshTokenStruct.RefreshToken, cfg)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["id"].(float64))

		if err := auth.ValidateRefreshToken(db, userID, refreshTokenStruct.RefreshToken); err != nil {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		user, err := auth.GetUserByID(db, userID)
		if err != nil {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
			return
		}

		newToken, newRefreshToken, err := auth.GenerateToken(user, cfg)
		if err != nil {
			http.Error(w, "Failed to generate new token", http.StatusInternalServerError)
			return
		}

		if err := auth.SaveRefreshToken(db, user.ID, newRefreshToken); err != nil {
			http.Error(w, "Failed to save new refresh token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": newToken, "refresh_token": newRefreshToken})
	}
}
