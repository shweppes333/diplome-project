package models

import "time"

type User struct {
    ID        int       `json:"id"`
    Login     string    `json:"login"`
    Password  string    `json:"password,omitempty"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

type Product struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Price       int    `json:"price"`
    Description string `json:"description"`
    Icon        string `json:"icon"`
    ImageURL    string `json:"image_url"`
}

type CartItem struct {
    ID          int    `json:"id"`
    UserID      int    `json:"user_id"`
    ProductID   int    `json:"product_id"`
    Quantity    int    `json:"quantity"`
    Price       int    `json:"price"`
    ProductName string `json:"product_name"`
    ProductIcon string `json:"product_icon"`
}

type RegisterRequest struct {
    Login    string `json:"login" binding:"required"`
    Password string `json:"password" binding:"required"`
    Email    string `json:"email" binding:"required"`
    Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
    Login    string `json:"login" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type AddToCartRequest struct {
    ProductID int `json:"product_id" binding:"required"`
    Quantity  int `json:"quantity"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}