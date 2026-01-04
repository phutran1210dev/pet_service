package middleware

import (
	"net/http"
	"pet-service/config"
	"pet-service/database"
	"pet-service/models"
	"pet-service/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserInfo struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"is_admin"`
	JTI       string `json:"jti"`
	TokenType string `json:"token_type"`
}

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.ErrCodeUnauthorized, "Authorization header required"))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.ErrCodeUnauthorized, "Invalid authorization format"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.SecretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.ErrCodeInvalidToken, utils.ErrorInvalidToken))
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.ErrCodeInvalidToken, utils.ErrorInvalidToken))
			c.Abort()
			return
		}

		// Check token type
		tokenType, _ := claims["token_type"].(string)
		if tokenType != "access_token" {
			c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.ErrCodeInvalidToken, utils.ErrorInvalidToken))
			c.Abort()
			return
		}

		// Check if token is blacklisted
		jti, _ := claims["jti"].(string)
		if isBlacklisted(jti) {
			c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.ErrCodeUnauthorized, utils.JTIInBlacklist))
			c.Abort()
			return
		}

		// Set user info in context
		userInfo := UserInfo{
			UserID:    claims["user_id"].(string),
			Username:  claims["username"].(string),
			FirstName: claims["first_name"].(string),
			LastName:  claims["last_name"].(string),
			Email:     claims["email"].(string),
			IsAdmin:   claims["is_admin"].(bool),
			JTI:       jti,
			TokenType: tokenType,
		}

		c.Set("current_user", userInfo)
		c.Next()
	}
}

// PermissionMiddleware checks if user has required permissions
func PermissionMiddleware(requiredPermissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, exists := c.Get("current_user")
		if !exists {
			c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.ErrCodeUnauthorized, "Unauthorized"))
			c.Abort()
			return
		}

		user := userInfo.(UserInfo)

		// Admin has all permissions
		if user.IsAdmin {
			c.Next()
			return
		}

		// Get user permissions from database
		db := database.GetDB()
		var userPermissions []string

		db.Table("permissions").
			Select("permissions.name").
			Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
			Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
			Where("user_roles.user_id = ? AND permissions.is_active = ?", user.UserID, true).
			Pluck("name", &userPermissions)

		if len(userPermissions) == 0 {
			c.JSON(http.StatusForbidden, utils.NewErrorResponse(utils.ErrCodePermissionDenied, utils.UserHasNoPermission))
			c.Abort()
			return
		}

		// Check if user has all required permissions
		permissionMap := make(map[string]bool)
		for _, p := range userPermissions {
			permissionMap[p] = true
		}

		for _, required := range requiredPermissions {
			if !permissionMap[required] {
				c.JSON(http.StatusForbidden, utils.NewErrorResponse(utils.ErrCodePermissionDenied, utils.PermissionDenied))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// isBlacklisted checks if token is in blacklist
func isBlacklisted(jti string) bool {
	db := database.GetDB()
	var count int64
	db.Model(&models.TokenBlacklist{}).Where("jti = ? AND is_active = ?", jti, true).Count(&count)
	return count > 0
}
