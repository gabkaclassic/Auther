package auth

import (
	"fmt"
	"time"

	"auther/configs"
	"auther/internal/db/models"
	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(user *models.User, cfg *configs.JWTConfig) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":  user.Login,
		"groups": user.Groups,
		"exp":    time.Now().Add(time.Duration(cfg.Expiration) * time.Second).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":  user.Login,
		"groups": user.Groups,
		"exp":    time.Now().Add(time.Duration(cfg.RefreshExpiration) * time.Second).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.RefreshSecret))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func parse(tokenString string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func ParseToken(tokenString string, cfg *configs.JWTConfig) (*jwt.Token, error) {
	return parse(tokenString, cfg.Secret)
}

func ParseRefreshToken(tokenString string, cfg *configs.JWTConfig) (*jwt.Token, error) {
	return parse(tokenString, cfg.RefreshSecret)
}
