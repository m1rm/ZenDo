package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"zendo/internal/database"
)

type Server struct {
	Port int

	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi("8090")
	NewServer := &Server{
		Port: port,

		db: database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
