package server

import (
	"log"
	"net/http"
	"wealthscope/backend/config"
)

func StartServer(cfg *config.Config) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", returnTestData)

    log.Printf("Starting server on port :%s", cfg.Port)
	err := http.ListenAndServe(":" + cfg.Port, mux)	
	log.Fatal(err)
}

func returnTestData(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s called", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json") // jzml
	w.Write([]byte(`{"message": "HELLO FROM BACKEND"}`))
}