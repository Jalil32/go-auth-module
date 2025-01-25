package server

import (
	"log"
	"net/http"
	"wealthscope/backend/config"

	"github.com/jmoiron/sqlx"
)

func StartServer(cfg *config.Config, db *sqlx.DB) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", returnTestData)

	// Test the database connection
	if err := db.Ping(); err != nil {
		return err
	} else {
		log.Println("Successfully connected")
	}

	log.Printf("Starting server on port :%s", cfg.Port)

	err := http.ListenAndServe(":"+cfg.Port, mux)

	return err
}

func returnTestData(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s called", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json") // jzml
	w.Write([]byte(`{"message": "HELLO FROM BACKEND"}`))
}
