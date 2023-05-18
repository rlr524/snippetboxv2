package models

import (
	"database/sql"
	"errors"
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
	// SQL statement that will be executed; use ? placeholders for values
	// not interpolation of variables to guard against injection attacks
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES (?, ?,
            UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use Exec() on the embedded connection pool to execute the statement. This returns a sql.Result
	// type, which contains basic information about what happened when the statement was executed.
	// Exec() compiles a prepared statement and stores it, then, in a next step, passes parameter values (?) to
	// the database where the DB executes the prepared statement using the parameters. Because the database is
	// getting the parameters after the statement is compiled, they're treated as pure data and can't change the intent
	// of the statement, so if a user inputs a statement intended as an injection attack, it will simply be
	// treated is any other query parameter, it can't actually be executed. This is required when preparing your
	// own sql statements as opposed to using methods provided by an ORM/ODM.
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
	// Statement that will be executed
	stmt := `SELECT id, title, content, created, expires FROM snippets
             WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Use QueryRow() method on the connection pool to execute the statement, passing in the untrusted
	// id variable as a value for the placeholder parameter. This returns a pointer to a sql.Row object
	// which holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct
	s := &Snippet{}

	// Use row.Scan() to copy the values from each field in sql.row to the corresponding field in the Snippet
	// struct. The arguments to row.Scan are *pointers* to the target for the copied data and the number of
	// arguments must be exactly the same as the number of columns returned by the statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a sql.ErrNoRows error. Use the errors.Is()
		// function to check for that error specifically, and return a custom ErrNoRecord error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// GetLatest returns a slice of instances of Snippet and a possible error
func (m *SnippetModel) GetLatest() ([]*Snippet, error) {
	return nil, nil
}
