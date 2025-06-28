package repositories

import (
	"database/sql"

	"github.com/alex-jienexa/labqueueueueue/models"
)

type QueueEntryRepository struct {
	db *sql.DB
}

func NewQueueEntryRepository(db *sql.DB) *QueueEntryRepository {
	return &QueueEntryRepository{db: db}
}

func (r *QueueEntryRepository) Create(queueEntry *models.QueueEntry) error {
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

// Получение позиции в очереди по его ID
func (r *QueueEntryRepository) GetByID(id int) (*models.QueueEntry, error) {
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
