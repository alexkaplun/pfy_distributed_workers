package storage

import (
	"database/sql"

	"github.com/alexkaplun/pfy_distributed_workers/model"
)

// For test task purposes I do not consider huge databases, and assume that all the URL are fetched by bulk
// I also considered reading the NEW records one by one and queuing them for workers, though not sure if in this
// particular case it would provide any efficiency boost. Simple approach will work good enough for the test task purpose
func (s *Storage) GetNewURLList() ([]*model.URL, error) {
	rows, err := s.db.Query(sqlSelectNewURLs)
	if err != nil {
		return nil, err
	}

	urls := make([]*model.URL, 0)
	for rows.Next() {
		item := model.URL{}
		var httpCodeNullable sql.NullInt32
		if err := rows.Scan(&item.Id, &item.Url, &item.Status, &httpCodeNullable); err != nil {
			return nil, err
		}

		if httpCodeNullable.Valid {
			item.HttpCode = httpCodeNullable.Int32
		} else {
			item.HttpCode = -1
		}

		urls = append(urls, &item)
	}
	return urls, nil
}

func (s *Storage) UpdateURLById(id int, status string, httpCode int) (bool, error) {
	res, err := s.db.Exec(sqlUpdateUrlById, status, httpCode, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	// return false if no rows were updated
	if rowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (s *Storage) UpdateStatusById(id int, status string) (bool, error) {
	res, err := s.db.Exec(sqlUpdateStatusById, status, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (s *Storage) GetURLById(id int) (*model.URL, error) {
	rows, err := s.db.Query(sqlSelectUrlById, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// if record does not exist for the provided id, return nil
	if !rows.Next() {
		return nil, nil
	}

	url := model.URL{}
	var httpCodeNullable sql.NullInt32
	if err := rows.Scan(&url.Id, &url.Url, &url.Status, &httpCodeNullable); err != nil {
		return nil, err
	}

	// handle nullable http code
	if httpCodeNullable.Valid {
		url.HttpCode = httpCodeNullable.Int32
	} else {
		url.HttpCode = -1
	}

	return &url, nil
}
