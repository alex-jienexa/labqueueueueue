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

// Получить все очереди
func (r *queueRepository) GetAll() ([]models.Queue, error) {
	query := `
		SELECT id, admin_id, title, is_active, starts_at, ends_at, conflict_resolution_method, created_at
		FROM queues
		ORDER BY id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.Queue
	for rows.Next() {
		var entry models.Queue
		err := rows.Scan(
			&entry.ID,
			&entry.AdminID,
			&entry.Title,
			&entry.IsActive,
			&entry.StartsAt,
			&entry.EndsAt,
			&entry.ResolutionMethod,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// Возврат активной очереди
func (r *queueRepository) GetActive() (*models.Queue, error) {
	// Обновляем активность очередей вне очереди
	r.ManageActive()

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

// Получает очередь по его ID
func (r *queueRepository) GetByID(id int) (*models.Queue, error) {
	query := `
		SELECT id, admin_id, title, is_active, starts_at, ends_at, conflict_resolution_method, created_at
		FROM queues
		WHERE id = $1
	`
	queue := &models.Queue{}
	err := r.db.QueryRow(query, id).Scan(
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
func (r *queueRepository) MoveAndPush(queueEntry *models.QueueEntry, position int) error {
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

// Вставляет уже созданный элемент насильно в позицию очереди.
// Если на данной позиции уже имеются элементы, то все элементы в нём отмечают как конфликтные
func (r *queueRepository) MoveForce(entry *models.QueueEntry, position int) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Откатываем в случае ошибки

	// Проверяем занята ли позиция перед изменением
	isBusy, err := r.IsPositionBusy(entry.QueueID, position)
	if err != nil {
		return fmt.Errorf("failed to check position: %w", err)
	}

	// Вставляем элемент
	err = tx.QueryRow(`
		UPDATE queue_entries
		SET position = $1
		WHERE id = $2
		RETURNING position
	`, position, entry.ID).Scan(&entry.Position)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("queue entry not found")
		} else {
			return fmt.Errorf("failed to update queue entry: %w", err)
		}
	}

	// Делаем все элементы в позиции конфликтными если до этого в нём были элементы
	if isBusy {
		_, err = tx.Exec(`
			UPDATE queue_entries
			SET is_conflict = TRUE
			WHERE queue_id = $1 AND position = $2
		`, entry.QueueID, entry.Position)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("queue entry not found")
			} else {
				return fmt.Errorf("failed to update queue entry: %w", err)
			}
		}
	}

	// Фиксируем изменения
	return tx.Commit()
}

// Перемещает уже существующий элемент очереди в первую свободную позицию дальше.
// Используется в разрешении конфликтных ситуаций
func (r *queueRepository) MoveToNextFree(entry *models.QueueEntry) error {
	if !entry.IsConflict {
		// Элемент не конфликтный
		// Todo: стоит ли что-то делать?
	}

	entries, err := r.GetEntries(entry.QueueID)
	for _, element := range entries {
		if element.Position > entry.Position {
			// Проверяем позицию следующую за element
			isBisy, err := r.IsPositionBusy(entry.ID, element.Position+1)
			if err != nil {
				return fmt.Errorf("failed to check position: %w", err)
			}
			if !isBisy {
				// Позиция свободна, туда вставляем элемент
				return r.MoveForce(entry, element.Position+1)
			}
		}
	}

	return err
}

// Проверяет, занята ли позиция position в очереди queueID.
func (r *queueRepository) IsPositionBusy(queueID int, position int) (bool, error) {
	var isBusy bool
	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM queue_entries
			WHERE queue_id = $1 AND position = $2
		)`, queueID, position).Scan(&isBusy)

	return isBusy, err
}

// Перемещает элемент в заданную позицию position очереди.
// Если в данной позиции уже имеется элемент, то он становится прямо за ним.
func (r *queueRepository) MoveAndFree(entry *models.QueueEntry, position int) error {
	if entry.Position == position {
		// Позиция не изменилась
		return nil
	} else {
		// Позиция изменилась
		// Проверяем занята ли позиция перед изменением
		if isBusy, err := r.IsPositionBusy(entry.QueueID, position); err != nil {
			return fmt.Errorf("failed to check position: %w", err)
		} else if isBusy {
			// Позиция занята
			// Пытаемся переместить элемент в следующую позицию
			return r.MoveAndFree(entry, position+1)
		} else {
			// Позиция свободна
			// Перемещаем элемент в свободную позицию
			return r.MoveForce(entry, position)
		}
	}
}

// Перепроверяет все элементы очереди и назначает IsActive = true для тех,
// у которых начинается время записи. Ровно и наоборот, если время записи уже
// закончилось, ставит IsActive как false
func (r *queueRepository) ManageActive() error {
	// Начнём транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Изменяем те очереди, чьё время настало
	if _, err := tx.Exec(`
		UPDATE queue_entries
		SET is_active = TRUE
		WHERE start_time <= NOW()
	`); err != nil {
		return fmt.Errorf("failed to update queue entries: %w", err)
	}

	// Изменяем те очереди, чьё время уже ушло
	if _, err := tx.Exec(`
		UPDATE queue_entries
		SET is_active = FALSE
		WHERE end_time <= NOW()
	`); err != nil {
		return fmt.Errorf("failed to update queue entries: %w", err)
	}

	return tx.Commit()
}
