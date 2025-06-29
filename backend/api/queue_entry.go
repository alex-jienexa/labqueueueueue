package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alex-jienexa/labqueueueueue/auth"
	"github.com/alex-jienexa/labqueueueueue/models"
	"github.com/alex-jienexa/labqueueueueue/repositories"
	"github.com/gin-gonic/gin"
)

func JoinQueue(c *gin.Context,
	queueRepo repositories.QueueRepository,
	studentRepo repositories.StudentRepository,
	queueEntryRepo repositories.QueueEntryRepository) {
	var input struct {
		Position int `json:"position"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный формат данных: " + err.Error()})
		return
	}

	queue_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный формат ID"})
		return
	}

	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Данные пользователя не найдены"})
		return
	}

	userClaims := claims.(*auth.Claims) // Приводим к типу Claims
	userID := userClaims.UserID

	student, err := studentRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении студента: " + err.Error()})
		return
	}
	queueEntry := models.QueueEntry{
		QueueID:   queue_id,
		StudentID: student.ID,
		Position:  input.Position,
		CreatedAt: time.Now(),
	}

	// Присоединяем студента в очередь
	err = queueEntryRepo.Create(&queueEntry)
	if err != nil {
		if err.Error() == "already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "Студент уже присоединился к очереди"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении студента в очередь: " + err.Error()})
		}
		return
	}
	// Моментально перемещаем его в зависимости от текущей позиции в очереди
	if queueEntry.Position == 0 {
		// Перемещаем его на первое свободное место
		if err = queueRepo.MoveToNextFree(&queueEntry); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при перемещении студента в очередь: " + err.Error()})
			return
		}
		log.Print("queueEntry: ", queueEntry)
		if err = queueEntryRepo.UpdateConflict(&queueEntry); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении позиции студента в очереди: " + err.Error()})
			return
		}
	} else {
		// Перемещаем в выборанное место
		if err := queueRepo.MoveForce(&queueEntry, queueEntry.Position); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при перемещении студента в очередь: " + err.Error()})
			return
		}
		if err = queueEntryRepo.UpdateConflict(&queueEntry); err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": "Ошибка при обновлении позиции студента в очереди: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, queueEntry)
}
