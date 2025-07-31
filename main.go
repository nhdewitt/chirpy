package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nhdewitt/chirpy/internal/http/handlers"
)

const (
	port         = "8080"
	filepathRoot = "."
)

func main() {
	godotenv.Load()

	// Establish DB Connection
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to DB: %s", err)
	}

	// API config
	apiCfg := handlers.NewAPIConfig(db)

	// Router
	mux := http.NewServeMux()
	apiCfg.RegisterRoutes(mux, filepathRoot)

	// Start server
	log.Printf("Server running at %s", port)
	addr := ":" + port
	log.Fatal(http.ListenAndServe(addr, mux))
}
