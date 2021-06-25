package manager

import (
	"context"

	"github.com/go-logr/logr"
	"golang.org/x/sync/errgroup"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Server holds details for a manager server.
type Server struct {
	mgr    manager.Manager
	cancel context.CancelFunc
	logger logr.Logger
}

// NewServer returns a new manager server.
func NewServer(logger logr.Logger, mgr manager.Manager) *Server {
	return &Server{
		logger: logger,
		mgr:    mgr,
	}
}

// Start starts the manager server.
func (s *Server) Start(ctx context.Context) error {
	shutdownContext, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := s.mgr.Start(shutdownContext); err != nil {
			s.logger.Error(err, "problem running manager")
			return err
		}
		return nil
	})
	return g.Wait()
}

// Stop signals the manager server to shutdown.
func (s *Server) Stop() {
	s.cancel()
}
