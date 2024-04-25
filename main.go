package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func main() {
	// Get the port number from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set.")
	}

	http.HandleFunc("/uid", uidHandler)

	// Start the HTTP server
	fmt.Printf("Server is running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func uidHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Generate a new UUID
	uid := uuid.New()
	fmt.Fprintf(w, "%s", uid.String())
}
