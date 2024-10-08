package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
	} `json:"database"`
}

func main() {
	config, err := readConfig("/config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	db := setupDatabase(config)
	defer db.Close()

	http.HandleFunc("/uid", func(w http.ResponseWriter, r *http.Request) {
		uidHandler(w, r, db)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, r, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
		log.Println("Set te listening port to 80")
	}

	log.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func readConfig(filePath string) (Config, error) {
	var config Config
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(bytes, &config)
	return config, err
}

func setupDatabase(cfg Config) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Ensure the database exists
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + cfg.Database.Dbname)
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	// Select the database
	_, err = db.Exec("USE " + cfg.Database.Dbname)
	if err != nil {
		log.Fatalf("Error selecting database: %v", err)
	}

	// Create the uids table if it does not exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS uids (
        id INT AUTO_INCREMENT PRIMARY KEY,
        uid VARCHAR(255) NOT NULL,
        timestamp DATETIME NOT NULL
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	return db
}

func uidHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}

	uid := uuid.New()
	_, err := db.Exec("INSERT INTO uids (uid, timestamp) VALUES (?, ?)", uid.String(), time.Now())
	if err != nil {
		http.Error(w, "Error saving UID", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	fmt.Fprintf(w, "%s", uid.String())
}

func healthHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if err := db.Ping(); err != nil {
		http.Error(w, "Database not accessible", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
