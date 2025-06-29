package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/alex-jienexa/labqueueueueue/auth"
	"github.com/alex-jienexa/labqueueueueue/models"
	"github.com/alex-jienexa/labqueueueueue/repositories"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context, repo repositories.StudentRepository) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Surname  string `json:"surname"`
		IsAdmin  bool   `json:"is_admin"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data: " + err.Error()})
		return
	}

	// Логирование для отладки
	log.Printf("Raw password during registration: %s", input.Password)

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Hashing error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	log.Printf("Generated hash: %s", string(hashedPassword))

	student := models.Student{
		Username: input.Username,
		Password: string(hashedPassword),
		Name:     input.Name,
		Surname:  input.Surname,
		IsAdmin:  input.IsAdmin,
	}

	log.Printf("Generated hash: %s", string(hashedPassword))

	if err := repo.Create(&student); err != nil {
		log.Printf("DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func Login(c *gin.Context, repo repositories.StudentRepository) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Добавляем логирование сырых данных
	// rawData, _ := c.GetRawData()
	// log.Printf("Raw login request: %s", string(rawData))

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Логируем полученные данные
	log.Printf("Login attempt for: %s", input.Username)

	student, err := repo.GetByUsername(input.Username)
	if err != nil {
		log.Printf("User not found: %s, error: %v", input.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин/пароль"})
		return
	}

	// Временная проверка хеширования
	testHash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	log.Printf("Test hash: %s", testHash)
	if string(testHash) == student.Password {
		log.Println("Manual hash check: OK")
	} else {
		log.Println("Manual hash check: FAIL")
	}

	// Логирование для отладки
	log.Printf("DB Hash: %s", student.Password)
	log.Printf("Input Pass: %s", input.Password)

	// Критически важный момент сравнения
	err = bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(input.Password))
	if err != nil {
		if err := bcrypt.CompareHashAndPassword(testHash, []byte(input.Password)); err == nil {
			log.Println("TEST: New hash works!")
		}
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Printf("Password mismatch for user %s", input.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин/пароль"})
			return
		}
		log.Printf("BCrypt error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	token, err := auth.GenerateToken(student.ID)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
