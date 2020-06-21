package storage

import (
	"strings"
)

func (s *Storage) InitDB() error {
	for _, query := range strings.Split(sqlInitDatabase, ";") {
		if len(strings.TrimSpace(query)) == 0 {
			continue
		}

		_, err := s.db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) MustHaveDB() error {
	// check that the 'urls' table exists in the DB
	rows, err := s.db.Query(sqlCheckUrlTableExists)
	if err != nil {
		return err
	}
	defer rows.Close()

	// in case the 'urls' table does not exists - run the init queries and fill the table with data
	if !rows.Next() {
		return s.InitDB()
	}
	return nil
}
