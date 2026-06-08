package main

import (
    "diplome-project/database"
    "diplome-project/handlers"
    "diplome-project/middleware"
    "log"

    "github.com/gin-gonic/gin"
)

func main() {
    // Инициализация базы данных
    database.InitDB()
    defer database.CloseDB()

    // Создание таблиц
    database.CreateTables()

    // Заполнение товаров
    database.SeedProducts()

    r := gin.Default()

    // Настройка CORS для фронтенда
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    // Публичные маршруты
    r.POST("/api/register", handlers.Register)
    r.POST("/api/login", handlers.Login)
    r.GET("/api/products", handlers.GetProducts)

    // Защищенные маршруты (требуют JWT)
    protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/cart", handlers.GetCart)
        protected.POST("/cart", handlers.AddToCart)
        protected.DELETE("/cart/:product_id", handlers.RemoveFromCart)
        protected.POST("/cart/clear", handlers.ClearCart)
    }

    // Запуск сервера с CORS
    log.Println("Server starting on :8080")
    r.Run(":8080")
}