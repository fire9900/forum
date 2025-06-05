package database

import (
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func NewSQLiteConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "../data.db")
	if err != nil {
		return nil, fmt.Errorf("(Forum) ошибка подключения к SQLite: %w", err)
	}
	db.SetMaxOpenConns(1000)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("(Forum) Не удалось проверить связь с SQLite: %w", err)
	}
	return db, nil
}
