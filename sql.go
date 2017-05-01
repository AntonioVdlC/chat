package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

// DBConfig holds the config values to connect to a PostgreSQL
// database
type DBConfig struct {
	Host string
	Port int64
	User string
	Password string
	Name string
	SSLMode string
}

// getDBConfig returns the config values to connect to the DB
// defaulting to a dev environment and reading values for the
// envarionment variables
func getDBConfig() DBConfig {
	config := DBConfig{
		Host: "localhost",
		Port: 5432,
		User: "dev",
		Password: "dev",
		Name: "chat_dev",
		SSLMode: "disable",
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		config.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.Port, _ = strconv.ParseInt(port, 10, 0)
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Password = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		config.Name = name
	}
	if sslmode := os.Getenv("DB_SSLMODE"); sslmode != "" {
		config.SSLMode = sslmode
	}

	return config
}

// initDB connects to the DB and creates the tables if they don't exist
func initDB() *sql.DB {
	config := getDBConfig()

	var (
		host = config.Host
		port = config.Port
		user = config.User
		password = config.Password
		dbname = config.Name
		sslmode = config.SSLMode
	)

	// Connect to DB
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	log.Println("Connected to DB.")

	// Create tables if not exists
	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id VARCHAR(255),
			user_name VARCHAR(255),
			user_avatar VARCHAR(255),
			type VARCHAR(255),
			content TEXT,
			date_post TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	log.Println("Tables created or already existing.")

	// All good!
	return db
}

// insertMessage inserts a single message into the database
// and returns either the id or an error
func insertMessage(db *sql.DB, msg Message) (string, error) {
	stmt := `
		INSERT INTO messages (
			user_id,
			user_name,
			user_avatar,
			type,
			content,
			date_post
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id string

	err := db.QueryRow(stmt, msg.UserID, msg.UserName, msg.UserAvatar, msg.Type, msg.Content, msg.Date).Scan(&id)
	
	if err != nil {
		return "", err
	}
	return id, nil
}

func selectPreviousMessage(db *sql.DB, userID string) (*sql.Rows, error) {
	stmt := `
		(SELECT *
		FROM messages
		WHERE type = 'message'
		LIMIT 10)
		
		UNION

		(SELECT *
		FROM messages
		WHERE type = 'message'
			AND user_id != $1
			AND date_post > (
				SELECT date_post
				FROM messages
				WHERE type = 'notice'
					AND content = 'logout'
					AND user_id = $1
				ORDER BY date_post DESC
				LIMIT 1
			)
		)
	`

	rows, err := db.Query(stmt, userID)
	if err != nil {
		return &sql.Rows{}, err
	}
	return rows, nil
}
