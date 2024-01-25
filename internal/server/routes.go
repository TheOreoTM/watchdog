package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/database", s.PrimaryDatabaseHandler)
	r.Get("/cache", s.CacheDatabaseHandler)
	r.Get("/", s.RootHandler)

	return r
}

func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Welcome to Watchdog. The VPS is online."

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error handling JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) PrimaryDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	err := s.db.Ping()
	if err != nil {
		resp["message"] = "Error connecting to database: " + err.Error()
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("error handling JSON marshal. Err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Set status to 500
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonResp)
		return
	}

	resp["message"] = "Primary database connection successful"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error handling JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) CacheDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	ctx := context.Background()

	status := s.redisClient.Ping(ctx)
	if status.Err() != nil {
		resp["message"] = "Error connecting to database " + status.Err().Error()
		jsonResp, err := json.Marshal(resp)

		if err != nil {
			log.Printf("error handling JSON marshal. Err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonResp)
		return
	}

	resp["message"] = "Cache database connection successful"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error handling JSON marshal. Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}
