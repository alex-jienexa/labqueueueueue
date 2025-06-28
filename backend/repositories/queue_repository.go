package repositories

import (
	"database/sql"

	"github.com/alex-jienexa/labqueueueueue/models"
)

type QueueRepository struct {
	db *sql.DB
}

func NewQueueRepository(db *sql.DB) *QueueRepository {
	return &QueueRepository{db: db}
}

// Создание новой очереди из модели models.Queue
func (r *QueueRepository) Create(queue *models.Queue) error {
	query := `
		INSERT INTO queues (admin_id, title, is_active, starts_at, ends_at, conflict_resolution_method)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		queue.AdminID,
		queue.Title,
		queue.IsActive,
		queue.StartsAt,
		queue.EndsAt,
		queue.ResolutionMethod,
	).Scan(&queue.ID)
	return err
}

// Возврат активной очереди
func (r *QueueRepository) GetActive() (*models.Queue, error) {
	query := `
		SELECT id, admin_id, title, is_active, starts_at, ends_at, conflict_resolution_method, created_at
		FROM queues
		WHERE is_active = TRUE
		LIMIT 1
	`
	queue := &models.Queue{}
	err := r.db.QueryRow(query).Scan(
		&queue.ID,
		&queue.AdminID,
		&queue.Title,
		&queue.IsActive,
		&queue.StartsAt,
		&queue.EndsAt,
		&queue.ResolutionMethod,
		&queue.CreatedAt,
	)
	return queue, err
}

// Получение всех записей в очередь QueueEntry для очереди Queue с ID `queueID`
func (r *QueueRepository) GetEntries(queueID int) ([]models.QueueEntry, error) {
	query := `
		SELECT id, queue_id, student_id, position, is_conflict, created_at
		FROM queue_entries
		WHERE queue_id = $1
		ORDER BY position
	`
	rows, err := r.db.Query(query, queueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.QueueEntry
	for rows.Next() {
		var entry models.QueueEntry
		err := rows.Scan(
			&entry.ID,
			&entry.QueueID,
			&entry.StudentID,
			&entry.Position,
			&entry.IsConflict,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
