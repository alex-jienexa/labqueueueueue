package api

import (
	"net/http"

	"github.com/alex-jienexa/labqueueueueue/models"
	"github.com/alex-jienexa/labqueueueueue/repositories"
	"github.com/gin-gonic/gin"
)

func CreateQueue(c *gin.Context, queueRepo repositories.QueueRepository, studRepo repositories.StudentRepository) {
	userID := c.MustGet("userId").(int)

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

	// Задать isActive в зависимости от времени начала и окончания
	if queue.StartsAt.Before(queue.EndsAt) {
		queue.IsActive = true
	} else {
		queue.IsActive = false
	}

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
