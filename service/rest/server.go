package rest

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codegangsta/negroni"

	. "fs/service/config"
	"fs/service/metrics"
)

// Server stores all needed fields for an API server
type Server struct {
	server   *http.Server
	shutdown chan struct{}
}

// Runner is an interface for the API server
type Runner interface {
	Run()
}

// NewServer initializes all server structures and starts shutdown listener
func NewServer() (*Server, error) {
	// basic negroni stuff init
	handler := negroni.New()
	router := NewRouter()
	router.Use(loggingMiddleware)
	router.Use(authMiddleware)
	handler.UseHandler(router)
	
	// set up of main server config structure
	server := &http.Server{
		Addr:           Setup.ServerIPAddress + ":" + Setup.ServerPort,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// server waits on it if interrupted
	shutdown := make(chan struct{})

	// Start shutdown handler as separate litener process
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint
		close(shutdown)
		Log.Info("Server shutdown requested")
		os.Exit(0)
	}()

	return &Server{server, shutdown}, nil
}

// RunServer starts the HTTP/HTTPS Sserver
func (s Server) RunServer() {
	Log.Info("Starting HTTP server on address: " + s.server.Addr)
	metrics.SetHealthy(metrics.ServiceHealthy)
	metrics.SetAlive(metrics.ServiceAlive)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		metrics.SetHealthy(metrics.ServiceError)
		metrics.SetAlive(metrics.ServiceDead)
		Log.Error("Error in ListenAndServe: " + err.Error())
	}

	// wait for shutdown signals
	<-s.shutdown
}
