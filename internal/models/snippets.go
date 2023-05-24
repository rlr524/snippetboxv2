package models

import (
	"database/sql"
	"errors"
	"fmt"
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

// We're not using an ORM, we're writing our own SQL statements and using the drivers directly via the database/sql
// package and the mysql drivers package. This is a bit more verbose than using an ORM, but our code is
// non-magical. The database/sql package works generally seamlessly with all popular SQL implementations
// so the DB functions are portable if we decide to switch from MySQL to PostgreSQL or another popular SQL DB.

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
	// Statement that will be executed
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP()
             ORDER BY id DESC LIMIT 10`

	// Use the Query() method on the connection pool to execute the statement.
	// This returns a sql.Rows result set.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Defer rows.Close() to ensure the sql.Rows() result set is always properly closed before the Latest()
	// method returns. The defer statement should come after checking for an error from the Query() method,
	// otherwise if Query() returns an error, the app will panic trying to close a nil result set.
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs
	snippets := []*Snippet{}
	fmt.Println("Snippets := []*Snippet{} ", snippets) // For testing

	// Use rows.Next() to iterate through the rows in the result set. This prepares the first (and then each
	// subsequent) row to be acted on by the rows.Scan() method. If iteration over all the rows completes, then
	// the result set automatically closes itself and frees up the underlying database connection.
	for rows.Next() {
		// Create a pointer to a now zeroed Snippet struct
		s := &Snippet{}
		fmt.Println("s := &Snippet{} ", s) // For testing
		// Use rows.Scan() to copy the values from each field in the row to the new Snippet object that
		// has been created. Again, the arguments to row.Scan() must be pointers to the target to which to
		// copy the data into, and the number of arguments must be exactly the same as the number of columns
		// returned by the sql statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		fmt.Println("s := &Snippet{} after scan -> ", s) // For testing
		// Append the object to the slice of snippets
		snippets = append(snippets, s)

	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error that was encountered
	// during the iteration. It's important to call this, don't assume that a successful iteration was
	// completed over the whole result set.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything is ok then return the Snippets slice
	return snippets, nil
}
