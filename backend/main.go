package main

import (
	"log"
	"os"

	"github.com/alex-jienexa/labqueueueueue/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация БД
	database.InitDB()
	defer database.DB.Close()

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

	// Запуск
	log.Printf("Server started on :%s", port)
	r.Run(":" + port)
}
