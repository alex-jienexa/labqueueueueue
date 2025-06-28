package models

// Conflict представляет конфликтную ситуацию, которая возникла при
// выбора места в очереди.
type Conflict struct {
	ID               int    `json:"id"`
	QueueID          int    `json:"queue_id"`
	StudentID        int    `json:"student_id"`
	Resolved         bool   `json:"resolved"`
	ResolutionMethod string `json:"resolution_method"` // "dice_roll", "manual", etc.
}
