package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(filename string) (*Storage, error) {
	handle, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: handle,
	}, nil
}

func (s *Storage) DB() *sql.DB {
	return s.db
}
