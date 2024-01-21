package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
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

func (s *Server) CacheDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       1,
	})

	ctx := context.Background()

	status := redisClient.Ping(ctx)
	if status.Err() != nil {
		resp["message"] = "Error connecting to database " + status.Err().Error()
		jsonResp, err := json.Marshal(resp)

		if err != nil {
			log.Fatalf("error handling JSON marshal. Err: %v", err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonResp)
		return
	}

	redisClient.Close()

	resp["message"] = "Cache database connection successful"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
