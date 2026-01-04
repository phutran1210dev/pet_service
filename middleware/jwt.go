package middleware

import (
	"pet-service/config"
	"pet-service/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken generates JWT access and refresh tokens
func GenerateToken(userID, username, firstName, lastName, email, jti string, isAdmin bool) (string, string, int64, error) {
	cfg := config.AppConfig

	// Access token expires in 1 day
	accessTokenExp := time.Now().Add(time.Duration(utils.AccessTokenExpire1Day) * time.Minute)
	accessClaims := jwt.MapClaims{
		"user_id":    userID,
		"username":   username,
		"first_name": firstName,
		"last_name":  lastName,
		"email":      email,
		"is_admin":   isAdmin,
		"jti":        jti,
		"token_type": "access_token",
		"exp":        accessTokenExp.Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return "", "", 0, err
	}

	// Refresh token expires in 1 year
	refreshTokenExp := time.Now().Add(365 * 24 * time.Hour)
	refreshClaims := jwt.MapClaims{
		"user_id":    userID,
		"username":   username,
		"first_name": firstName,
		"last_name":  lastName,
		"email":      email,
		"is_admin":   isAdmin,
		"jti":        jti,
		"token_type": "refresh_token",
		"exp":        refreshTokenExp.Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return "", "", 0, err
	}

	return accessTokenString, refreshTokenString, accessTokenExp.Unix(), nil
}

// GetCurrentUser gets current user from context
func GetCurrentUser(c *gin.Context) (UserInfo, bool) {
	userInfo, exists := c.Get("current_user")
	if !exists {
		return UserInfo{}, false
	}
	return userInfo.(UserInfo), true
}
