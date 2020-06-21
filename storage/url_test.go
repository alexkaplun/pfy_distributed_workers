package storage

import (
	"fmt"
	"log"
	"testing"

	"github.com/alexkaplun/pfy_distributed_workers/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TEST_DB_FILENAME = "test.db"

var storage *Storage

func TestStorageGetURLById(t *testing.T) {
	cases := map[string]struct {
		id       int
		expected *model.URL
	}{
		"missing id": {
			id:       999,
			expected: nil,
		},
		"existing id NEW": {
			id: 3,
			expected: &model.URL{
				Id: 3, Url: "http://go.org", Status: "NEW", HttpCode: -1,
			},
		},
		"existing id DONE": {
			id: 5,
			expected: &model.URL{
				Id: 5, Url: "https://google.com", Status: "DONE", HttpCode: 200,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			url, err := storage.GetURLById(tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, url)
		})
	}
}

// assume to test number of records with the "NEW" status
func TestStorageGetNewURLList(t *testing.T) {
	totalNew, err := storage.GetNewURLList()
	require.NoError(t, err)
	assert.Equal(t, 4, len(totalNew))
}

func TestStorageUpdateURLById(t *testing.T) {
	testUrl := &model.URL{
		Id: 777, Url: "https://test.com", Status: model.StatusProcessing, HttpCode: 200,
	}
	expectedUrl := &model.URL{
		Id: 777, Url: "https://test.com", Status: model.StatusDone, HttpCode: 200,
	}
	// create a record we are going to update
	sqlInsert := fmt.Sprintf(`INSERT INTO urls VALUES (%v, "%v", "%v", %v)`,
		testUrl.Id, testUrl.Url, testUrl.Status, testUrl.HttpCode)
	_, err := storage.db.Exec(sqlInsert)
	assert.NoError(t, err)

	defer func() {
		// clean up the record
		sqlDelete := fmt.Sprintf("DELETE FROM urls WHERE id = %v", testUrl.Id)
		_, err = storage.db.Exec(sqlDelete)
		require.NoError(t, err)
	}()

	// test update for non-existing id
	res, err := storage.UpdateURLById(999, "DONE", 200)
	assert.NoError(t, err)
	assert.False(t, res)

	// update the newly created record
	res, err = storage.UpdateURLById(testUrl.Id, string(expectedUrl.Status), int(expectedUrl.HttpCode))
	require.NoError(t, err)

	// check that the record now looks as needed
	url, err := storage.GetURLById(testUrl.Id)
	require.NoError(t, err)
	assert.Equal(t, expectedUrl, url)
}

func TestStorageUpdateStatusById(t *testing.T) {
	testUrl := &model.URL{
		Id: 777, Url: "https://test.com", Status: model.StatusNew, HttpCode: -1,
	}
	expectedUrl := &model.URL{
		Id: 777, Url: "https://test.com", Status: model.StatusError, HttpCode: -1,
	}
	// create a record we are going to update

	sqlInsert := fmt.Sprintf(`INSERT INTO urls VALUES (%v, "%v", "%v", %v)`,
		testUrl.Id, testUrl.Url, testUrl.Status, testUrl.HttpCode)
	_, err := storage.db.Exec(sqlInsert)
	assert.NoError(t, err)

	defer func() {
		// clean up the record
		sqlDelete := fmt.Sprintf("DELETE FROM urls WHERE id = %v", testUrl.Id)
		_, err = storage.db.Exec(sqlDelete)
		require.NoError(t, err)
	}()

	// test update for non-existing id
	res, err := storage.UpdateStatusById(999, "ERROR")
	assert.NoError(t, err)
	assert.False(t, res)

	// update the newly created record
	res, err = storage.UpdateURLById(testUrl.Id, string(expectedUrl.Status), int(expectedUrl.HttpCode))
	require.NoError(t, err)

	// check that the record now looks as needed
	url, err := storage.GetURLById(testUrl.Id)
	require.NoError(t, err)
	assert.Equal(t, expectedUrl, url)
}

func init() {
	s, err := New(TEST_DB_FILENAME)
	if err != nil {
		log.Fatal("error initializing db handle:", err)
	}

	if err := s.MustHaveDB(); err != nil {
		log.Fatal("error creating the database:", err)
	}
	storage = s
}
