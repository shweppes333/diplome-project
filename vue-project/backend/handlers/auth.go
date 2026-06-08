package handlers

import (
	"database/sql"
	"diplome-project/database"
	"diplome-project/models"
	"net/http"
	"regexp"
	"time"
	"log" 

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

func generateToken(userID int, login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"login":   login,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func validateLogin(login string) bool {
	return len(login) >= 6
}

func validatePassword(password string) bool {
	return len(password) >= 6
}

func validateEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}


func Register(c *gin.Context) {
    log.Println("=== НАЧАЛО ОБРАБОТКИ РЕГИСТРАЦИИ ===")
    
   
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Printf(" Ошибка парсинга JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные запроса"})
        return
    }
    
    log.Printf(" Получены данные: Login=%s, Email=%s, Name=%s, PasswordLen=%d", 
        req.Login, req.Email, req.Name, len(req.Password))

   
    if len(req.Login) < 6 {
        log.Printf(" Логин слишком короткий: %s (длина: %d)", req.Login, len(req.Login))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Логин должен содержать минимум 6 символов"})
        return
    }
    log.Println(" Логин прошел валидацию")

    
    if len(req.Password) < 6 {
        log.Printf(" Пароль слишком короткий: длина %d", len(req.Password))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Пароль должен содержать минимум 6 символов"})
        return
    }
    log.Println(" Пароль прошел валидацию")


    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(req.Email) {
        log.Printf(" Неверный формат email: %s", req.Email)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный email адрес"})
        return
    }
    log.Println(" Email прошел валидацию")

    // 5. Проверка существования пользователя
    var existingLogin string
    err := database.DB.QueryRow(
        "SELECT login FROM users WHERE login = $1 OR email = $2", 
        req.Login, req.Email,
    ).Scan(&existingLogin)
    
    if err == nil {
        log.Printf(" Пользователь уже существует: login=%s, email=%s", req.Login, req.Email)
        c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким логином или email уже существует"})
        return
    }
    if err != sql.ErrNoRows {
        log.Printf(" Ошибка при проверке пользователя: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера при проверке пользователя"})
        return
    }
    log.Println("✅ Пользователь не найден, можно создавать")

    // 6. Хеширование пароля
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf(" Ошибка хеширования пароля: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке пароля"})
        return
    }
    log.Println("✅ Пароль успешно захэширован")

    // 7. Вставка в БД
    var userID int
    err = database.DB.QueryRow(`
        INSERT INTO users (login, password, email, name)
        VALUES ($1, $2, $3, $4)
        RETURNING id`,
        req.Login, string(hashedPassword), req.Email, req.Name,
    ).Scan(&userID)
    
    if err != nil {
        log.Printf("❌ Ошибка при сохранении пользователя в БД: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении пользователя"})
        return
    }

    log.Printf(" ПОЛЬЗОВАТЕЛЬ УСПЕШНО СОЗДАН! ID=%d", userID)
    log.Println("=== КОНЕЦ РЕГИСТРАЦИИ ===")
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Регистрация успешна",
        "user_id": userID,
    })
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Поиск пользователя в БД
	var user models.User
	var hashedPassword string
	err := database.DB.QueryRow(`
        SELECT id, login, password, email, name
        FROM users
        WHERE login = $1`,
		req.Login).Scan(&user.ID, &user.Login, &hashedPassword, &user.Email, &user.Name)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	// Генерация JWT токена
	token, err := generateToken(user.ID, user.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  user,
	})
}
