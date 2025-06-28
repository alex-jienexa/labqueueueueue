package models

import "time"

// QueueEntry предоставляет позицию некоторого студента StudentID в
// очереди QueueID.
type QueueEntry struct {
	ID         int       `json:"id"`
	QueueID    int       `json:"queue_id"`
	StudentID  int       `json:"student_id"`
	Position   int       `json:"position"`
	IsConflict bool      `json:"is_conflict"`
	CreatedAt  time.Time `json:"created_at"`
}
