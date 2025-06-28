package models

import "time"

// Queue предоставляет объект очереди, который содержит информацию о
// очереди. Тот, кто создаёт очередь автоматически становится его
// администратором. Запись в очередь может быть только в промежутке
// со временем начала и окончания очереди.
// Параметр ResolutionMethod определяет, как в очереди будет решаться конфликт.
type Queue struct {
	ID               int       `json:"id"`
	AdminID          int       `json:"admin_id"`
	Title            string    `json:"title"`
	IsActive         bool      `json:"is_active"`
	StartsAt         time.Time `json:"starts_at"`
	EndsAt           time.Time `json:"ends_at"`
	ResolutionMethod string    `json:"resolution_method"` // "move_after", "first_free", "to_end"
	CreatedAt        time.Time `json:"created_at"`
}
