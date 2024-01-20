package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/database", s.PrimaryDatabaseHandler)

	return r
}

func (s *Server) PrimaryDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	db, err := sql.Open("postgres", os.Getenv("PRIMARY_DB_URI"))
	if err != nil {
		resp["message"] = "Error connecting to database" + err.Error()
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("error handling JSON marshal. Err: %v", err)
		}
		// Set status to 500
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonResp)
		return
	}
	if db != nil {
		db.Close()
	}

	resp["message"] = "Primary database connection successful"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
