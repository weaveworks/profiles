package interrupt

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/sync/errgroup"
)

const timeout = 15 * time.Second

// Service defines a service which can start and stop itself.
type Service interface {
	Start(ctx context.Context) error
	Stop()
}

// Handler collects all services it handles interrupts for.
type Handler struct {
	logger   logr.Logger
	services []Service
}

// NewInterruptHandler return a new interrupt handler for the given services.
func NewInterruptHandler(logger logr.Logger, services ...Service) *Handler {
	logger = logger.WithName("handler")
	return &Handler{
		logger:   logger,
		services: services,
	}
}

// ListenAndGracefulShutdown monitors services for interrupts. If received, it starts a graceful shutdown of all the services it
// monitors. If the deadline for the graceful shutdown is reached, it will terminate all processes with os.Exit(1).
func (h *Handler) ListenAndGracefulShutdown() error {
	g, ctx := errgroup.WithContext(context.Background())
	for _, s := range h.services {
		s := s
		g.Go(func() error { return s.Start(ctx) })
	}

	g.Go(func() error {
		interruptChannel := make(chan os.Signal, 2)
		signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

		<-interruptChannel
		done := make(chan struct{})
		// start the timer for the shutdown sequence
		go func() {
			select {
			case <-done:
				return
			case <-time.After(timeout):
				h.logger.Error(errors.New("timeout"), "graceful shutdown timed out... forcing shutdown")
				os.Exit(1)
			}
		}()
		h.logger.Info("received shutdown signal... gracefully terminating servers...")
		for _, s := range h.services {
			s.Stop()
		}
		h.logger.Info("all done. Goodbye.")
		done <- struct{}{}
		return nil
	})
	return g.Wait()
}
