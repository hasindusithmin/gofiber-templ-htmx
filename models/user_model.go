package models

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

const bcryptCost = 8

func CreateUser(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcryptCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (email, password, username) VALUES (?, ?, ?);`
	_, err = db.Exec(query, user.Email, string(hashedPassword), user.Username)
	return err
}

func GetUserById(id string) (User, error) {
	query := `SELECT id, email, password, username FROM users WHERE id = ?`
	return scanUserRow(db.QueryRow(query, id))
}

func CheckEmail(email string) (User, error) {
	query := `SELECT id, email, password, username FROM users WHERE email = ?`
	return scanUserRow(db.QueryRow(query, email))
}

// Helper function to scan a user from a sql.Row
func scanUserRow(row *sql.Row) (User, error) {
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil // or return custom error if needed
		}
		return User{}, err
	}
	return user, nil
}
