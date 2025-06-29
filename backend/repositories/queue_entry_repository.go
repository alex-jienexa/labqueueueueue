package repositories

import (
	"database/sql"
	"errors"

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
		INSERT INTO queue_entries (queue_id, student_id, position, is_conflict, created_at)
		VALUES ($1, $2, $3, $4, $5)
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
		FROM queue_entries
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

func (r *queueEntryRepository) UpdateConflict(queueEntry *models.QueueEntry) error {
	// Выбираем все элементы в позици текущей
	if queueEntry.Position > 0 {
		query := `
			SELECT id, queue_id, student_id, position, is_conflict, created_at
			FROM queue_entries
			WHERE position = $1 AND queue_id = $2
		`
		rows, err := r.db.Query(
			query,
			queueEntry.Position,
			queueEntry.QueueID,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return err
		}
		defer rows.Close()

		var countOnPosition int
		err = r.db.QueryRow(`SELECT COUNT(id) FROM queue_entries WHERE position = $1 AND queue_id = $2`, queueEntry.Position, queueEntry.QueueID).Scan(&countOnPosition)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return err
		}
		isMultiple := countOnPosition > 1
		queueEntry.IsConflict = isMultiple

		for rows.Next() {
			var entry models.QueueEntry
			err = rows.Scan(
				&entry.ID,
				&entry.QueueID,
				&entry.StudentID,
				&entry.Position,
				&entry.IsConflict,
				&entry.CreatedAt,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil
				}
				return err
			}
			_, err = r.db.Exec(`UPDATE queue_entries SET is_conflict = $1 WHERE id = $2`, isMultiple, entry.ID)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := errors.New("position is not set")
	return err
}
