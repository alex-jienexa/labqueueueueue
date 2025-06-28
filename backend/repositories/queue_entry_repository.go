package repositories

import (
	"database/sql"

	"github.com/alex-jienexa/labqueueueueue/models"
)

type queueEntryRepository struct {
	db *sql.DB
}

func NewQueueEntryRepository(db *sql.DB) QueueEntryRepository {
	return &queueEntryRepository{db: db}
}

func (r *queueEntryRepository) Create(queueEntry *models.QueueEntry) error {
	query := `
		INSERT INTO queues (queue_id, student_id, position, is_conflict, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		queueEntry.QueueID,
		queueEntry.StudentID,
		queueEntry.Position,
		queueEntry.IsConflict,
		queueEntry.CreatedAt,
	).Scan(&queueEntry.ID)
	return err
}

// Добавляет элемент в очередь. Позиция элемента зависит от значения queueEntry.Position.
// Если queueEntry.Position == 0, то элемент добавляется в первый свободный слот очереди
func (r *queueEntryRepository) AddToQueue(queueEntry *models.QueueEntry, queueRepo queueRepository) error {
	if queueEntry.Position == 0 {
		// Позиция не указана, добавляем в конец очереди
		r.Create(queueEntry)
		err := queueRepo.MoveToNextFree(queueEntry)
		if err != nil {
			return err
		}
	}
	if queueEntry.Position > 0 {
		// Позиция указана, добавляем в указанную позицию
		r.Create(queueEntry)
		err := queueRepo.MoveForce(queueEntry, queueEntry.Position)
		if err != nil {
			return err
		}
	}

	return nil
}

// Получение позиции в очереди по его ID
func (r *queueEntryRepository) GetByID(id int) (*models.QueueEntry, error) {
	query := `
		SELECT id, queue_id, student_id, position, is_conflict, created_at
		FROM queues
		WHERE id = $1
	`
	queueEntry := &models.QueueEntry{}
	err := r.db.QueryRow(query, id).Scan(
		&queueEntry.ID,
		&queueEntry.QueueID,
		&queueEntry.StudentID,
		&queueEntry.Position,
		&queueEntry.IsConflict,
		&queueEntry.CreatedAt,
	)

	return queueEntry, err
}
