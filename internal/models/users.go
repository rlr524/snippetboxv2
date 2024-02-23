package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
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
func (m *UserModel) Insert(username, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// Use the Exec() method to insert the user details and hashed password into the users table.
	_, err = m.DB.Exec(stmt, username, email, string(hashedPassword))
	if err != nil {
		// If this returns an error, we use the errors.As() function to check whether the error has the type
		// *mysql.MySQLError. If it does, the error will be assigned to the mySQLError variable. We can then check
		// whether the error relates to our users_uc_email key by checking if the error code equals 1062 and
		// the contents of the error message string. If it does, we return an ErrDuplicateEmail error.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
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
