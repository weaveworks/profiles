package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	gruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/weaveworks/profiles/pkg/protos"
)

const timeout = 10 * time.Second

// Server contains details for the gateway server.
type Server struct {
	logger   logr.Logger
	server   *http.Server
	apiAddr  string
	grpcAddr string
	cancel   context.CancelFunc
}

// NewGatewayServer creates a new grpc-gateway server.
func NewGatewayServer(logger logr.Logger, apiAddr string, grpcAddr string) *Server {
	logger = logger.WithName("gateway-server")
	return &Server{
		logger:   logger,
		apiAddr:  apiAddr,
		grpcAddr: grpcAddr,
	}
}

// Start starts the grpc-gateway server using from Endpoint.
func (s *Server) Start() error {
	// setup grpc-gateway to connect to the grpc server
	mux := gruntime.NewServeMux()
	gopts := []grpc.DialOption{grpc.WithInsecure()}
	if err := protos.RegisterProfilesServiceHandlerFromEndpoint(context.Background(), mux, s.grpcAddr, gopts); err != nil {
		s.logger.Error(err, "failed to register service handler from endpoint")
		return err
	}

	s.logger.Info(fmt.Sprintf("starting profiles grpc-gateway server at %s", s.apiAddr))
	server := &http.Server{Addr: s.apiAddr, Handler: mux}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		// ignore server is closing error because the server receives that on graceful shutdown.
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(err, "unable to start profiles api server")
			return err
		}
		return nil
	})
	s.server = server
	return g.Wait()
}

// Stop does a graceful shutdown of the server using a timeout of 10 seconds.
func (s *Server) Stop() {
	serverTimeoutContext, timeout := context.WithTimeout(context.Background(), timeout)
	defer timeout()
	if err := s.server.Shutdown(serverTimeoutContext); err != nil {
		s.logger.Error(err, "Failed to gracefully shutdown server... terminating.")
	}
	s.cancel()
	s.logger.Info("server stopped")
}
