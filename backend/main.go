package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	_ "github.com/mattn/go-sqlite3"
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

// NEED to implement CORS as middleware for security
func main() {
	db := initDB()
	defer db.Close()

	createTable(db)

	mux := http.NewServeMux()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)(mux)

	mux.HandleFunc("/api/auth/google", googleAuthHandler)

	fmt.Println("server is running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
