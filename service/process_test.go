package service

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexkaplun/pfy_distributed_workers/model"

	"github.com/alexkaplun/pfy_distributed_workers/storage"
)

var s *Service

const TEST_DB_FILENAME = "test.db"

func TestProcessUrl_Success(t *testing.T) {
	urlId := 777
	testUrl := &model.URL{
		Id:       urlId,
		Url:      "https://proxify.io",
		Status:   "NEW",
		HttpCode: -1,
	}

	// create a record we are going to update
	sqlInsert := fmt.Sprintf(`INSERT INTO urls VALUES (%v, "%v", "%v", %v)`,
		testUrl.Id, testUrl.Url, testUrl.Status, testUrl.HttpCode)
	_, err := s.storage.DB().Exec(sqlInsert)
	assert.NoError(t, err)

	defer func() {
		// clean up the record
		sqlDelete := fmt.Sprintf("DELETE FROM urls WHERE id = %v", testUrl.Id)
		_, err = s.storage.DB().Exec(sqlDelete)
		require.NoError(t, err)
	}()

	err = s.processUrl(testUrl)
	assert.NoError(t, err)

	// assert the url status
	url, err := s.storage.GetURLById(urlId)
	require.NoError(t, err)
	assert.Equal(t, model.StatusDone, url.Status)
}

func TestProcessUrl_Missing(t *testing.T) {
	testUrl := &model.URL{
		Id:       777,
		Url:      "https://proxify.io",
		Status:   "NEW",
		HttpCode: -1,
	}

	// attempt to process url missing from the DB
	err := s.processUrl(testUrl)
	assert.Error(t, err)
}

func TestProcessUrl_WrongUrl(t *testing.T) {
	urlId := 777
	testUrl := &model.URL{
		Id:       urlId,
		Url:      "https://bad url",
		Status:   "NEW",
		HttpCode: -1,
	}

	// create a record we are going to update
	sqlInsert := fmt.Sprintf(`INSERT INTO urls VALUES (%v, "%v", "%v", %v)`,
		testUrl.Id, testUrl.Url, testUrl.Status, testUrl.HttpCode)
	_, err := s.storage.DB().Exec(sqlInsert)
	assert.NoError(t, err)

	defer func() {
		// clean up the record
		sqlDelete := fmt.Sprintf("DELETE FROM urls WHERE id = %v", testUrl.Id)
		_, err = s.storage.DB().Exec(sqlDelete)
		require.NoError(t, err)
	}()

	// must be error due to http get error
	err = s.processUrl(testUrl)
	assert.Error(t, err)

	// assert the url status as ERROR
	url, err := s.storage.GetURLById(urlId)
	require.NoError(t, err)
	assert.Equal(t, model.StatusError, url.Status)
}

func init() {
	db, err := storage.New(TEST_DB_FILENAME)
	if err != nil {
		log.Fatal("error initializing db handle:", err)
	}

	if err := db.MustHaveDB(); err != nil {
		log.Fatal("error creating the database:", err)
	}

	s = New(db, 1)
}
