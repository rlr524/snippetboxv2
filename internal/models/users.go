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
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         int8
}

type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to the "users" table.
func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// Use the Exec() method to insert the user details and hashed password into the users table.
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	// Why not create a method to check the db for the email vs depending on the MySQL error number, which
	// MySQL could change and tightly couples this method to MySQL? Because that method introduces a race
	// condition to the application. If two users try to sign up with the same email at exactly the same time,
	// both submissions will pass the validation check but only one INSERT statement will succeed and the other
	// will violate the UNIQUE constraint set on email in the database and get a 500 error. Neither case is
	// optimal and even though the race condition is extremely unlikely and fairly benign, this is all probably a
	// good reason why we should be using an ORM instead of rolling our own SQL.
	// TODO: At some point, determine how to optimize this. Use an ORM?
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
	// Retrieve the id and hashed password associated with the given email.
	// If no matching email exists, we return the ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check whether hashed_password and the plain-text password provided match.
	// If they don't, return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

// Exists checks if a user exists given a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
