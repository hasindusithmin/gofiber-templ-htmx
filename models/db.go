package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	once sync.Once
)

// init initializes the database and performs migrations once.
func init() {
	if err := InitDB(); err != nil {
		log.Fatalf("🔥 Failed to initialize DB: %v", err)
	}
}

// InitDB ensures the DB connection and runs migrations once.
func InitDB() error {
	var err error
	once.Do(func() {
		err = connect()
		if err != nil {
			return
		}
		err = makeMigrations()
	})
	return err
}

// connect establishes the DB connection.
func connect() error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error opening DB connection: %w", err)
	}

	if err = conn.Ping(); err != nil {
		return fmt.Errorf("error pinging DB: %w", err)
	}

	db = conn
	log.Println("🚀 Connected Successfully to the Database")
	return nil
}

// makeMigrations creates necessary tables if they don't exist.
func makeMigrations() error {
	userTable := `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		username VARCHAR(64) NOT NULL
	);`

	todoTable := `CREATE TABLE IF NOT EXISTS todos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		created_by INT NOT NULL,
		title VARCHAR(64) NOT NULL,
		description VARCHAR(255),
		status BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (created_by) REFERENCES users(id)
	);`

	if _, err := db.Exec(userTable); err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	if _, err := db.Exec(todoTable); err != nil {
		return fmt.Errorf("error creating todos table: %w", err)
	}
	return nil
}
