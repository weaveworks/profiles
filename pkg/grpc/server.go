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
	cancel   context.CancelFunc
}

// NewGRPCServer returns a new grpc server.
func NewGRPCServer(logger logr.Logger, catalog *catalog.Catalog, grpcAddr string) *Server {
	logger = logger.WithName("grpc")
	return &Server{
		logger:   logger,
		grpcAddr: grpcAddr,
		catalog:  catalog,
	}
}

// Start starts the grpc server on the given address and saves the servers state for later termination.
func (g *Server) Start() error {
	// setup grpc server details
	grpcLis, err := net.Listen("tcp", g.grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on address %s: %v", g.grpcAddr, err)
	}
	grpcSrv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	g.server = grpcSrv
	reflection.Register(grpcSrv)

	// create the catalog grpc server
	catalogGrpcServer := api.NewCatalogAPI(g.catalog, g.logger.WithName("api"))
	protos.RegisterProfilesServiceServer(grpcSrv, catalogGrpcServer)
	// serve grpc apis
	g.logger.Info(fmt.Sprintf("starting profiles grpc server at %s", g.grpcAddr))
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel
	e, _ := errgroup.WithContext(ctx)
	e.Go(func() error {
		if err := grpcSrv.Serve(grpcLis); err != nil {
			g.logger.Error(err, "unable to start grpc api server")
			return err
		}
		return nil
	})
	return e.Wait()
}

// Stop does a graceful shutdown of the grpc server.
func (g *Server) Stop() {
	g.server.GracefulStop()
	g.cancel()
	g.logger.Info("server stopped")
}
