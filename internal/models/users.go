package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	UserName       string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         int8
}

type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to the "users" table.
func (m *UserModel) Insert(userName, email, password string) error {
	return nil
}

// Authenticate verifies whether a user exists with the provided email address and password and
// returns the relevant user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Exists checks if a user exists given a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
