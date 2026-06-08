package handlers

import (
   "diplome-project/database"
    "diplome-project/models"
    "net/http"

    "github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
    rows, err := database.DB.Query(`
        SELECT id, name, price, description, icon, COALESCE(image_url, '')
        FROM products
        ORDER BY id
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения товаров"})
        return
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var p models.Product
        err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Icon, &p.ImageURL)
        if err != nil {
            continue
        }
        products = append(products, p)
    }

    c.JSON(http.StatusOK, products)
}