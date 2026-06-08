package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные claims"})
            c.Abort()
            return
        }

        userIDFloat, ok := claims["user_id"].(float64)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный user_id"})
            c.Abort()
            return
        }

        c.Set("user_id", int(userIDFloat))
        c.Set("login", claims["login"])
        c.Next()
    }
}