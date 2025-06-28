package repositories

import (
	"database/sql"
	"fmt"

	"github.com/alex-jienexa/labqueueueueue/models"
)

type queueRepository struct {
	db *sql.DB
}

func NewQueueRepository(db *sql.DB) QueueRepository {
	return &queueRepository{db: db}
}

// Создание новой очереди из модели models.Queue
func (r *queueRepository) Create(queue *models.Queue) error {
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
func (r *queueRepository) GetActive() (*models.Queue, error) {
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
func (r *queueRepository) GetEntries(queueID int) ([]models.QueueEntry, error) {
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

// Вставляет уже существующий элемент очереди в позицию position очереди.
// При этом, меняет позиции ВСЕХ элементов после position на +1
func (r *queueRepository) ForceMove(queueEntry *models.QueueEntry, position int) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Откатываем в случае ошибки

	// Сдвиг всего на position+1
	_, err = tx.Exec(
		`UPDATE queue_entries
		SET position = position + 1
		WHERE queue_id = $1 AND position >= $2
	`, queueEntry.QueueID, position)
	if err != nil {
		return fmt.Errorf("failed to shift queue entries: %w", err)
	}

	// Перемещение элемента
	err = tx.QueryRow(`
		UPDATE queue_entries
		SET position = $1
		WHERE id = $2
		RETURNING position
	`, position, queueEntry.ID).Scan(&queueEntry.Position)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("queue entry not found")
		} else {
			return fmt.Errorf("failed to update queue entry: %w", err)
		}
	}

	// Фиксируем изменения
	return tx.Commit()
}
