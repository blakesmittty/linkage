package main

import (
	"database/sql"
	"fmt"
	"log"
)

//user logs in using oauth, high score is retrieved and sent to the client

// database initialization
func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./highscores.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func createTable(db *sql.DB) {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS highscores (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        google_user_id TEXT NOT NULL,
        high_score INTEGER NOT NULL
    );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

func main() {
	db := initDB()
	defer db.Close()

	createTable(db)

	fmt.Printf("helloWorld!")
}
