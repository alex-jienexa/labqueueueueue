package main

import (
	"log"
	"os"

	"github.com/alex-jienexa/labqueueueueue/api"
	"github.com/alex-jienexa/labqueueueueue/database"
	"github.com/alex-jienexa/labqueueueueue/repositories"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация БД
	db := database.InitDB()
	defer db.Close()

	// Загружаем все репозитории
	studentRepo := repositories.NewStudentRepository(db)

	// Читаем порт из .env (по умолчанию 8080)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Настройка сервера
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Auth
	r.POST("/register", func(c *gin.Context) { api.Register(c, studentRepo) })
	r.POST("/login", func(c *gin.Context) { api.Login(c, studentRepo) })

	// Защита роутинга = требуют проверки JWT-токена
	// authGroup := r.Group("/")
	// authGroup.Use(api.JWTAuthMiddleware()) {
	// Вставить защищённые методы
	// }

	// Запуск
	log.Printf("Server started on :%s", port)
	r.Run(":" + port)
}
