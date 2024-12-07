package testdb

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	// Импортируем PostgreSQL-драйвер
	_ "github.com/lib/pq"
)

// Для конфигурации тестовой базы
const (
	envDSN   = "TEST_DB_DSN"
	postgres = "postgres"
)

// Имена таблиц
const (
	UsersTable   = "users"
	IncomesTable = "incomes"
)

// DB хранит соединение с тестовой базой данных
var DB *sql.DB

// SetupTestDB инициализирует подключение к тестовой базе
func SetupTestDB() error {
	dsn := os.Getenv(envDSN)
	if dsn == "" {
		return fmt.Errorf("environment variable %s is not set. Please set it to connect to the test database",
			envDSN)
	}

	db, err := sql.Open(postgres, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to test db: %w", err)
	}

	DB = db
	return nil
}

// CloseTestDB закрывает соединение с тестовой базой
func CloseTestDB() {
	if DB != nil {
		_ = DB.Close()
		DB = nil
	}
}

// TruncateTables очищает таблицы
func TruncateTables(db *sql.DB, tables ...string) error {
	if len(tables) == 0 {
		return nil
	}

	query := fmt.Sprintf("truncate table %s restart identity cascade;",
		strings.Join(tables, ", "))

	_, err := db.Exec(query)
	return err
}

// UserParams содержит параметры для создания тестового пользователя.
type UserParams struct {
	ID    int
	Email string
	Name  string
}

// SeedUser добавляет тестового пользователя.
// Указание параметров необязательно.
func (p *UserParams) SeedUser(db *sql.DB) error {
	if p.ID == 0 {
		p.ID = 1
	}
	if p.Email == "" {
		p.Email = "test@test.com"
	}
	if p.Name == "" {
		p.Name = "test tester"
	}

	_, err := db.Exec(`INSERT INTO users (id, email, name) VALUES ($1, $2, $3)`, p.ID, p.Email, p.Name)
	if err != nil {
		return err
	}
	return nil
}
