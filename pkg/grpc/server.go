package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/go-logr/logr"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/weaveworks/profiles/pkg/api"
	"github.com/weaveworks/profiles/pkg/catalog"
	"github.com/weaveworks/profiles/pkg/protos"
)

// Server contains details for the grpc server.
type Server struct {
	logger   logr.Logger
	grpcAddr string
	server   *grpc.Server
	catalog  *catalog.Catalog
}

// NewServer returns a new grpc server.
func NewServer(logger logr.Logger, catalog *catalog.Catalog, grpcAddr string) *Server {
	logger = logger.WithName("grpc")
	return &Server{
		logger:   logger,
		grpcAddr: grpcAddr,
		catalog:  catalog,
	}
}

// Start starts the grpc server on the given address and saves the servers state for later termination.
func (s *Server) Start(ctx context.Context) error {
	// setup grpc server details
	grpcLis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on address %s: %v", s.grpcAddr, err)
	}
	grpcSrv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	s.server = grpcSrv
	reflection.Register(grpcSrv)

	// create the catalog grpc server
	catalogGrpcServer := api.NewCatalogAPI(s.catalog, s.logger.WithName("api"))
	protos.RegisterProfilesServiceServer(grpcSrv, catalogGrpcServer)
	// serve grpc apis
	s.logger.Info(fmt.Sprintf("starting profiles grpc server at %s", s.grpcAddr))
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := grpcSrv.Serve(grpcLis); err != nil {
			s.logger.Error(err, "unable to start grpc api server")
			return err
		}
		return nil
	})
	return g.Wait()
}

// Stop does a graceful shutdown of the grpc server.
func (s *Server) Stop() {
	s.server.GracefulStop()
	s.logger.Info("server stopped")
}
