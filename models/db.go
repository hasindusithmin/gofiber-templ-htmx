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

func init() {
	MakeMigrations()
}

func getConnection() *sql.DB {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("MYSQL_USERNAME"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"),
		)

		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("🔥 Failed to connect to the database: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("🔥 Database ping failed: %v", err)
		}

		log.Println("🚀 Connected successfully to the database")
	})

	return db
}

func MakeMigrations() {
	conn := getConnection()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			username VARCHAR(64) NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS todos (
			id INT AUTO_INCREMENT PRIMARY KEY,
			created_by INT NOT NULL,
			title VARCHAR(64) NOT NULL,
			description VARCHAR(255),
			status BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (created_by) REFERENCES users(id)
		);`,
	}

	for _, stmt := range statements {
		if _, err := conn.Exec(stmt); err != nil {
			log.Fatalf("❌ Failed to execute migration: %v", err)
		}
	}
}
