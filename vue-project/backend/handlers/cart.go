package handlers

import (
    "database/sql"
    "diplome-project/database"
    "diplome-project/models"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

func GetCart(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
        return
    }

    rows, err := database.DB.Query(`
        SELECT c.id, c.product_id, c.quantity, c.price, p.name, p.icon
        FROM cart c
        JOIN products p ON c.product_id = p.id
        WHERE c.user_id = $1
        ORDER BY c.added_at DESC`,
        userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения корзины"})
        return
    }
    defer rows.Close()

    var cartItems []models.CartItem
    for rows.Next() {
        var item models.CartItem
        err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Price, &item.ProductName, &item.ProductIcon)
        if err != nil {
            continue
        }
        item.UserID = userID.(int)
        cartItems = append(cartItems, item)
    }

    c.JSON(http.StatusOK, cartItems)
}

func AddToCart(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
        return
    }

    var req models.AddToCartRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
        return
    }

    if req.Quantity <= 0 {
        req.Quantity = 1
    }

    // Получаем информацию о товаре
    var product models.Product
    err := database.DB.QueryRow(`
        SELECT id, name, price, icon
        FROM products
        WHERE id = $1`,
        req.ProductID).Scan(&product.ID, &product.Name, &product.Price, &product.Icon)
    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
        return
    }

    // Проверяем, есть ли уже товар в корзине
    var existingID int
    var existingQuantity int
    err = database.DB.QueryRow(`
        SELECT id, quantity FROM cart
        WHERE user_id = $1 AND product_id = $2`,
        userID, req.ProductID).Scan(&existingID, &existingQuantity)

    if err == sql.ErrNoRows {
        // Добавляем новый товар
        _, err = database.DB.Exec(`
            INSERT INTO cart (user_id, product_id, quantity, price, product_name, product_icon)
            VALUES ($1, $2, $3, $4, $5, $6)`,
            userID, req.ProductID, req.Quantity, product.Price, product.Name, product.Icon)
    } else if err == nil {
        // Обновляем количество
        _, err = database.DB.Exec(`
            UPDATE cart
            SET quantity = quantity + $1
            WHERE user_id = $2 AND product_id = $3`,
            req.Quantity, userID, req.ProductID)
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления в корзину"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Товар добавлен в корзину"})
}

func RemoveFromCart(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
        return
    }

    productID, err := strconv.Atoi(c.Param("product_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID товара"})
        return
    }

 
    result, err := database.DB.Exec(`
        DELETE FROM cart
        WHERE user_id = $1 AND product_id = $2`,
        userID, productID)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления товара"})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден в корзине"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Товар удален из корзины"})
}

func ClearCart(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
        return
    }

    _, err := database.DB.Exec(`
        DELETE FROM cart
        WHERE user_id = $1`,
        userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка очистки корзины"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Корзина очищена"})
}