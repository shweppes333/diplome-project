package database

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
    _ "golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "root"
    dbname   = "postgres"
)

func InitDB() {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    var err error
    DB, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal("Ошибка подключения к БД:", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal("Не удалось подключиться к БД:", err)
    }

    log.Println("Успешное подключение к PostgreSQL!")
}

func CloseDB() {
    if DB != nil {
        DB.Close()
    }
}

func CreateTables() {
    // Таблица пользователей
    userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        login VARCHAR(100) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
    
    _, err := DB.Exec(userTable)
    if err != nil {
        log.Fatal("Ошибка создания таблицы users:", err)
    }

    // Таблица товаров
    productsTable := `
    CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        price INTEGER NOT NULL,
        description TEXT,
        icon VARCHAR(100),
        image_url TEXT
    )`
    
    _, err = DB.Exec(productsTable)
    if err != nil {
        log.Fatal("Ошибка создания таблицы products:", err)
    }

    // Таблица корзины
    cartTable := `
    CREATE TABLE IF NOT EXISTS cart (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
        quantity INTEGER NOT NULL DEFAULT 1,
        price INTEGER NOT NULL,
        product_name VARCHAR(255) NOT NULL,
        product_icon VARCHAR(100),
        added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(user_id, product_id)
    )`
    
    _, err = DB.Exec(cartTable)
    if err != nil {
        log.Fatal("Ошибка создания таблицы cart:", err)
    }

    log.Println("Таблицы успешно созданы")
}

func SeedProducts() {
    var count int
    err := DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
    if err != nil {
        log.Println("Ошибка проверки товаров:", err)
        return
    }

    if count > 0 {
        log.Println("Товары уже существуют в БД")
        return
    }

    products := []struct {
        name        string
        price       int
        description string
        icon        string
    }{
        {"Керамическая кружка", 1250, "Ручная лепка, глазурь", "bi bi-cup-hot-fill"},
        {"Вязаный плед", 3450, "Хлопок, мягкий", "bi bi-border-width"},
        {"Деревянная ложка", 780, "Ручная резьба", "bi bi-egg-fried"},
        {"Свеча ручной работы", 540, "Аромат ванили", "bi bi-candle"},
        {"Сумка из кожи", 699, "Натуральная кожа, ручная работа", "bi bi-bag"},
        {"Цепочка-браслет", 299, "Серебро, ручное плетение", "bi bi-gem"},
    }

    for _, p := range products {
        _, err := DB.Exec(`
            INSERT INTO products (name, price, description, icon)
            VALUES ($1, $2, $3, $4)`,
            p.name, p.price, p.description, p.icon)
        if err != nil {
            log.Printf("Ошибка вставки товара %s: %v", p.name, err)
        }
    }

    log.Println("Товары успешно добавлены в БД")
}