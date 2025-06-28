package models

// StudentQueueAccess определяет уровни доступа студента в очереди.
// Для администратора очереди все уровни доступа разрешены.
type StudentQueueAccess struct {
	StudentID int  `json:"student_id"`
	QueueID   int  `json:"queue_id"`
	CanView   bool `json:"can_view"`
	CanJoin   bool `json:"can_join"`
}
