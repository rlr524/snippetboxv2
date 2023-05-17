package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

// Remember that using a receiver function is the same as declaring a method. These functions below
// could be declared inside the SnippetModel struct, but doing dependency injection this way
// essentially makes the functions static, meaning we don't need to use them with an instance of SnippetModel. This
// is also a good paradigm for testing as it makes it easy to create an interface and mock it for unit testing.

// Insert takes in a title, some content, and an expiration number of days and returns an id and possibly an error
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// SQL statement that will be executed
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES (?, ?,
            UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use Exec() on the embedded connection pool to execute the statement. This returns a sql.Result
	// type, which contains basic information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result to get the ID of the newly inserted record.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type of int64, so it's converted to an int before returning.
	return int(id), nil
}

// Get takes in an id and returns an instance of Snippet and a possible error
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// GetLatest returns a slice of instances of Snippet and a possible error
func (m *SnippetModel) GetLatest() ([]*Snippet, error) {
	return nil, nil
}
