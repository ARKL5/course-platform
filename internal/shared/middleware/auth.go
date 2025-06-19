package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 密钥 - 需要与service中的保持一致
var jwtSecret = []byte("course-platform-secret-key-2024")

// JWTClaims JWT声明结构体
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "请提供认证Token",
				"code":  "MISSING_TOKEN",
			})
			c.Abort()
			return
		}

		// 检查Token格式 (Bearer <token>)
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token格式错误，请使用 'Bearer <token>' 格式",
				"code":  "INVALID_TOKEN_FORMAT",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// 解析和验证Token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			var errorMessage string
			var errorCode string

			switch {
			case err == jwt.ErrSignatureInvalid:
				errorMessage = "Token签名无效"
				errorCode = "INVALID_SIGNATURE"
			case err == jwt.ErrTokenExpired:
				errorMessage = "Token已过期，请重新登录"
				errorCode = "TOKEN_EXPIRED"
			case err == jwt.ErrTokenMalformed:
				errorMessage = "Token格式错误"
				errorCode = "MALFORMED_TOKEN"
			default:
				errorMessage = "Token验证失败"
				errorCode = "INVALID_TOKEN"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errorMessage,
				"code":  errorCode,
			})
			c.Abort()
			return
		}

		// 检查Token是否有效
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token无效",
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// 提取用户信息
		if claims, ok := token.Claims.(*JWTClaims); ok {
			// 将用户信息存储到上下文中
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)

			// 继续处理请求
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token claims解析失败",
				"code":  "INVALID_CLAIMS",
			})
			c.Abort()
			return
		}
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求Token）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有Token，继续处理但不设置用户信息
			c.Next()
			return
		}

		// 有Token则尝试解析
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
			tokenString := tokenParts[1]

			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtSecret, nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(*JWTClaims); ok {
					c.Set("userID", claims.UserID)
					c.Set("username", claims.Username)
				}
			}
		}

		c.Next()
	}
}

// GetUserFromContext 从上下文中获取用户信息
func GetUserFromContext(c *gin.Context) (userID uint, username string, exists bool) {
	userIDInterface, exists1 := c.Get("userID")
	usernameInterface, exists2 := c.Get("username")

	if !exists1 || !exists2 {
		return 0, "", false
	}

	userID, ok1 := userIDInterface.(uint)
	username, ok2 := usernameInterface.(string)

	if !ok1 || !ok2 {
		return 0, "", false
	}

	return userID, username, true
}
