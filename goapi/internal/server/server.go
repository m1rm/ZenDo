package server

import (
	"net/http"
	"time"
	"zendo/internal/router"
	"zendo/internal/database"
)

func NewServer() *http.Server {

	rr := router.NewRouter(database.New())

	server := &http.Server{
		Addr:         ":8090",
		Handler:      rr.GetMux(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
