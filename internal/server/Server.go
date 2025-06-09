package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lalo64GG/api-gateway/internal/config"
	"github.com/Lalo64GG/api-gateway/internal/router"
	"github.com/Lalo64GG/api-gateway/internal/services"
)

type Server struct {
	server *http.Server
	cfg *config.Config
	handler http.Handler
}

func New(cfg *config.Config) *Server {
	serviceRegistry, err := services.NewRegistry(cfg)

	if err != nil {
		log.Fatalf("Error initializing service registry: %v", err)
	}

	handler := router.SetupRoutes(cfg, serviceRegistry)

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: handler,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	return &Server{
		server: srv,
		cfg: cfg,
		handler: handler,
	}
}


func (s *Server) Start() error {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)

	go func() {
		log.Printf("Starting API Gateway on %s", s.cfg.ServerAddr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case <-stop: 
		log.Println("Shutting down server...")
	case err:= <- errChan:
		log.Printf("Server error: %v", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
		return err
	}

	log.Println("Server gracefully stopped")
	return nil
}