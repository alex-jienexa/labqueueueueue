package repositories

import "github.com/alex-jienexa/labqueueueueue/models"

type StudentRepository interface {
	Create(student *models.Student) error
	GetByID(id int) (*models.Student, error)
	GetByUsername(username string) (*models.Student, error)
	GetByUsernameNoPassword(username string) (*models.Student, error)
}

type QueueRepository interface {
	Create(queue *models.Queue) error
	GetActive() (*models.Queue, error)
	GetEntries(queueID int) ([]models.QueueEntry, error)
}

type QueueEntryRepository interface {
	Create(queueEntry *models.QueueEntry) error
	GetByID(id int) (*models.QueueEntry, error)
}
