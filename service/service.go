package service

import (
	"context"

	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"github.com/nshumoogum/fantasy-football/api"
	"github.com/nshumoogum/fantasy-football/config"
	"github.com/pkg/errors"
)

//go:generate moq -out mock/server.go -pkg mock . HTTPServer
//go:generate moq -out mock/closer.go -pkg mock . Closer

// HTTPServer defines the required methods from the HTTP server
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// Closer defines the required methods for a closable resource
type Closer interface {
	Close(ctx context.Context) error
}

// Service contains all the configs, server and clients to run the Dataset API
type Service struct {
	api    *api.FantasyFootballAPI
	config *config.Configuration
	server HTTPServer
}

// New creates a new service
func New(cfg *config.Configuration) *Service {
	svc := &Service{
		api:    &api.FantasyFootballAPI{FPLURL: cfg.FPLURL},
		config: cfg,
	}

	return svc
}

// Run the service
func (svc *Service) Run(ctx context.Context, svcErrors chan error) (err error) {
	log.Event(ctx, "starting service", log.INFO)

	// Get HTTP router and server with middleware
	router := mux.NewRouter()
	svc.api = api.NewFantasyFootballAPI(ctx, svc.api.FPLURL, router)

	server := server.New(svc.config.BindAddr, router)

	// Disable this here to allow main to manage graceful shutdown of the entire app.
	server.HandleOSSignals = false

	svc.server = server

	// Run the http server in a new go-routine
	go func() {
		if err := svc.server.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.config.GracefulShutdownTimeout
	log.Event(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout}, log.INFO)
	shutdownContext, cancel := context.WithTimeout(ctx, timeout)
	hasShutdownError := false

	// Gracefully shutdown the application closing any open resources.
	go func() {
		defer cancel()

		// stop any incoming requests
		if err := svc.server.Shutdown(shutdownContext); err != nil {
			log.Event(shutdownContext, "failed to shutdown http server", log.Error(err), log.ERROR)
			hasShutdownError = true
		}

	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-shutdownContext.Done()

	// timeout expired
	if shutdownContext.Err() == context.DeadlineExceeded {
		log.Event(shutdownContext, "shutdown timed out", log.ERROR, log.Error(shutdownContext.Err()))
		return shutdownContext.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Event(shutdownContext, "failed to shutdown gracefully ", log.ERROR, log.Error(err))
		return err
	}

	log.Event(shutdownContext, "graceful shutdown was successful", log.INFO)
	return nil
}
