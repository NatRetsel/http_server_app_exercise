package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/natretsel/http_server_app_exercise/internal/database"
)

type apiConfig struct {
	filePathRoot string
	db           *database.Queries
	jwtSecret    string
}

func main() {
	// load env variables
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}

	filePathRoot := os.Getenv("FILEPATH_ROOT")
	if filePathRoot == "" {
		log.Fatal("FILEPATH_ROOT must be set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error trying to establish connection to DB %v: %v", dbURL, err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	// Initialize apiConfig struct
	cfg := apiConfig{
		filePathRoot: filePathRoot,
		db:           dbQueries,
		jwtSecret:    jwtSecret,
	}

	// initialize mux and httpserver
	mux := http.NewServeMux()

	// add routes and handler functions
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}
	log.Printf("Serving on: http://localhost:%s/app/\n", port)
	log.Fatal(srv.ListenAndServe())
}
