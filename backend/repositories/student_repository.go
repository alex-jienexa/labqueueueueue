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
		INSERT INTO students (name, surname, username, password_hash, is_admin)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		student.Name,
		student.Surname,
		student.Username,
		student.Password,
		student.IsAdmin,
	).Scan(&student.ID)

	return err
}

// Получение пользователя по ID
func (r *StudentRepository) GetByID(id int) (*models.Student, error) {
	query := `
		SELECT id, name, surname, username, is_admin
		FROM students
		WHERE id = $1
	`
	student := &models.Student{}
	err := r.db.QueryRow(query, id).Scan(
		&student.ID,
		&student.Name,
		&student.Surname,
		&student.Username,
		&student.IsAdmin,
	)
	return student, err
}

// Получение пользователя по username для авторизации
func (r *StudentRepository) GetByUsername(username string) (*models.Student, error) {
	query := `
		SELECT id, name, surname, username, password_hash, is_admin 
		FROM students 
		WHERE username = $1
	`
	student := &models.Student{}
	err := r.db.QueryRow(query, username).Scan(
		&student.ID,
		&student.Name,
		&student.Surname,
		&student.Password,
		&student.IsAdmin,
	)
	return student, err
}

// Получение пользователя по username но без пароля (для безопаности)
func (r *StudentRepository) GetByUsernameNoPassword(username string) (*models.Student, error) {
	query := `
		SELECT id, name, surname, username, is_admin 
		FROM students 
		WHERE username = $1
	`
	student := &models.Student{}
	err := r.db.QueryRow(query, username).Scan(
		&student.ID,
		&student.Name,
		&student.Surname,
		&student.IsAdmin,
	)
	return student, err
}
