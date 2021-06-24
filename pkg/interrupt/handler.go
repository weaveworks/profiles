package interrupt

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-logr/logr"

	"github.com/weaveworks/profiles/pkg/gateway"
	"github.com/weaveworks/profiles/pkg/grpc"
)

const timeout = 15 * time.Second

// Handler collects all services it handles interrupts for.
type Handler struct {
	logger            logr.Logger
	grpcServer        *grpc.Server
	gatewayServer     *gateway.Server
	managerCancelFunc context.CancelFunc
}

// NewInterruptHandler return a new interrupt handler for the given services.
func NewInterruptHandler(logger logr.Logger, grpcServer *grpc.Server, gatewayServer *gateway.Server, managerCancelFunc context.CancelFunc) *Handler {
	logger = logger.WithName("handler")
	return &Handler{
		logger:            logger,
		grpcServer:        grpcServer,
		gatewayServer:     gatewayServer,
		managerCancelFunc: managerCancelFunc,
	}
}

// HandleInterrupts monitors for interrupts. If received, it starts a graceful shutdown of all the services it
// monitors. If the deadline for the graceful shutdown is reached, it will terminate all processes with os.Exit(1).
func (h *Handler) HandleInterrupts() {
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
	// shutdown the manager
	h.managerCancelFunc()
	// shutdown grpc-gateway server
	h.gatewayServer.Stop()
	// shutdown grpc server
	h.grpcServer.Stop()
	h.logger.Info("all done. Goodbye.")
	done <- struct{}{}
}
