package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alex-jienexa/labqueueueueue/auth"
	"github.com/alex-jienexa/labqueueueueue/models"
	"github.com/alex-jienexa/labqueueueueue/repositories"
	"github.com/gin-gonic/gin"
)

func CreateQueue(c *gin.Context, queueRepo repositories.QueueRepository, studRepo repositories.StudentRepository) {
	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Данные пользователя не найдены"})
		return
	}

	userClaims := claims.(*auth.Claims) // Приводим к типу Claims
	userID := userClaims.UserID

	student, err := studRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	} else if !student.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Только админы могут создавать очереди",
		})
		return
	}

	// Получаем входные данные
	var queue models.Queue
	if err := c.ShouldBindJSON(&queue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	queue.AdminID = userID

	if queue.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Название очереди не может быть пустым",
		})
	}

	// ResolutionMethod может принимать только значения: `"move_after"`, `"first_free"`, `"to_end"`
	switch queue.ResolutionMethod {
	case "move_after", "first_free", "to_end":
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный метод разрешения очереди, принимает только move_after, first_free, to_end",
		})
		return
	}

	currentTime := time.Now()
	queue.IsActive = currentTime.After(queue.StartsAt) && currentTime.Before(queue.EndsAt)

	if err := queueRepo.Create(&queue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при создании очереди",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Очередь успешно создана",
		"queue":   queue,
	})
}

func GetActiveQueue(c *gin.Context, queueRepo repositories.QueueRepository) {
	queue, err := queueRepo.GetActive()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении активной очереди",
		})
		return
	} else if queue == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "На данный момент нет активных очередей",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"queue": queue,
	})
}

func GetQueueByID(c *gin.Context, queueRepo repositories.QueueRepository) {
	// Получить id из из параметров в URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный формат ID"})
		return
	}

	queue, err := queueRepo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Очередь не найдена"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении очереди"})
		}
		return
	}

	c.JSON(http.StatusOK, queue)
}

func GetAllQueues(c *gin.Context, queueRepo repositories.QueueRepository) {
	if queues, err := queueRepo.GetAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении всех очередей",
		})
	} else {
		c.JSON(http.StatusOK, queues)
	}
}
