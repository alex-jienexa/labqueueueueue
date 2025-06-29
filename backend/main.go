package main

import (
	"log"
	"os"

	"github.com/alex-jienexa/labqueueueueue/api"
	"github.com/alex-jienexa/labqueueueueue/database"
	"github.com/alex-jienexa/labqueueueueue/middleware"
	"github.com/alex-jienexa/labqueueueueue/repositories"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация БД
	db := database.InitDB()
	defer db.Close()

	// Загружаем все репозитории
	studentRepo := repositories.NewStudentRepository(db)
	queueRepo := repositories.NewQueueRepository(db)

	// Читаем порт из .env (по умолчанию 8080)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Настройка сервера
	r := gin.Default()

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Auth
	r.POST("/register", func(c *gin.Context) { api.Register(c, studentRepo) })
	r.POST("/login", func(c *gin.Context) { api.Login(c, studentRepo) })

	// Защита роутинга = требуют проверки JWT-токена
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("/rping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "rpong :3"}) })
		authGroup.POST("/queues", func(c *gin.Context) { api.CreateQueue(c, queueRepo, studentRepo) })
		authGroup.GET("/queues/active", func(c *gin.Context) { api.GetActiveQueue(c, queueRepo) })
		authGroup.GET("/queues/:id", func(c *gin.Context) { api.GetQueueByID(c, queueRepo) })
	}

	// Запуск
	log.Printf("Server started on :%s", port)
	r.Run(":" + port)
}
