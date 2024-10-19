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
        username TEXT PRIMARY KEY,
        high_score INTEGER NOT NULL DEFAULT 0
    );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

func insertNewUser(db *sql.DB, username string) {
	insertSQL := `INSERT INTO highscores (username) VALUES (?)`
	_, err := db.Exec(insertSQL, username)
	if err != nil {
		log.Fatalf("Error inserting new user: %v", err)
	}
}

func updateHighScore(db *sql.DB, username string, newScore int) {
	updateSQL := `
    UPDATE highscores
    SET high_score = ?
    WHERE username = ? AND high_score < ?`

	_, err := db.Exec(updateSQL, newScore, username, newScore)
	if err != nil {
		log.Fatalf("Error updating high score: %v", err)
	}
}

func getUserHighScore(db *sql.DB, username string) (int, error) {
	var highScore int
	query := `SELECT high_score FROM highscores WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&highScore)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // User not found, return 0 as the default score
		}
		return 0, err // Other errors
	}
	return highScore, nil
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

	//mux.HandleFunc("/ws", websocketHandler)

	mux.HandleFunc("/api/auth/google", googleAuthHandler)

	fmt.Println("server is running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
