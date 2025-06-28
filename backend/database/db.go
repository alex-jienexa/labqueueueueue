package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Драйвер PostgreSQL

	// Менеджеры миграций

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *sql.DB

func InitDB() {
	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Читаем переменные окружения
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Формируем строку подключения
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	psqlLink := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)
	Migrate("file://migrations", psqlLink)

	// Подключаемся к БД
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to PostgreSQL!")
}

// Выполняет миграции из директории `migrationsPath` для pgSQL-базы
// данных с ссылкой подключения `pgsqlInfo`.
func Migrate(migrationsPath string, pgsqlInfo string) {
	m, err := migrate.New(
		migrationsPath,
		pgsqlInfo)

	if err != nil {
		log.Fatalf("Error creating migration instance: %v", err)
	}
	errors := m.Up()
	if errors != nil && errors != migrate.ErrNoChange {
		log.Fatal(errors)
	} else if errors == migrate.ErrNoChange {
		log.Print("No migrations to apply.")
	} else {
		log.Print("Applied migrations!")
	}
}
