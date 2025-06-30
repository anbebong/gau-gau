package middleware

import (
	"net/http"
	"strings"
	"time"

	// Thay đổi thành import path thật sự của logutil

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte
var jwtExpire time.Duration

// Nên gọi hàm này ở api.Start để inject config
func InitJWT(secret string, expire time.Duration) {
	jwtSecret = []byte(secret)
	jwtExpire = expire
}

func JWTAuthMiddleware(handler gin.HandlerFunc, requireAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		if requireAdmin {
			role, _ := claims["role"].(string)
			if role != "admin" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
				return
			}
		}
		userID, _ := claims["user_id"].(string)
		username, _ := claims["username"].(string)
		role, _ := claims["role"].(string)
		// logutil.APIDebug("JWTAuthMiddleware: user_id = %v, username = %v, role = %v", userID, username, role)
		c.Set("user_id", userID)    // user_id (string) dùng cho mọi truy vấn
		c.Set("username", username) // username chỉ để hiển thị
		c.Set("role", role)
		handler(c)
	}
}

// Middleware cho group: chỉ kiểm tra JWT, không cần handler
func JWTAuthMiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		userID, _ := claims["user_id"].(string)
		username, _ := claims["username"].(string)
		role, _ := claims["role"].(string)
		// logutil.APIDebug("JWTAuthMiddlewareFunc: user_id = %v, username = %v, role = %v", userID, username, role)
		c.Set("user_id", userID)    // user_id (string) dùng cho mọi truy vấn
		c.Set("username", username) // username chỉ để hiển thị
		c.Set("role", role)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}
