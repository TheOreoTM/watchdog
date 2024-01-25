package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	db          *sql.DB
	redisClient *redis.Client
	port        int
}

func NewServer() *http.Server {
	db, err := sql.Open("postgres", os.Getenv("PRIMARY_DB_URI"))
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       1,
	})

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	s := &Server{
		db:          db,
		redisClient: redisClient,
		port:        port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.RegisterRoutes(),
	}

	return server
}
