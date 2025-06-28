package repositories

import (
	"database/sql"

	"github.com/alex-jienexa/labqueueueueue/models"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// Добавление нового пользователя в базу данных
func (r *StudentRepository) Create(student *models.Student) error {
	query := `
		INSERT INTO students (name, surname, password_hash, is_admin)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		student.Name,
		student.Surname,
		student.Password,
		student.IsAdmin,
	).Scan(&student.ID)

	return err
}
