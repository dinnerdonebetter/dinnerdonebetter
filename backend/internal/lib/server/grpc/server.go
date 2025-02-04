package grpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"

	"google.golang.org/grpc"
)

const (
	serviceName = "grpc_server"
)

type (
	Config struct {
		Port uint16 `env:"PORT" json:"port"`
	}

	Server struct {
		logger                   logging.Logger
		config                   *Config
		grpcServer               *grpc.Server
		unaryServerInterceptors  []grpc.UnaryServerInterceptor
		streamServerInterceptors []grpc.StreamServerInterceptor
	}

	RegistrationFunc func(*grpc.Server)
)

func NewGRPCServer(
	cfg *Config,
	logger logging.Logger,
	unaryServerInterceptors []grpc.UnaryServerInterceptor,
	streamServerInterceptors []grpc.StreamServerInterceptor,
	registrationFunctions ...RegistrationFunc,
) (*Server, error) {
	if cfg == nil {
		return nil, internalerrors.NilConfigError("grpc server")
	}

	opts := []grpc.ServerOption{}
	for _, interceptor := range unaryServerInterceptors {
		opts = append(opts, grpc.UnaryInterceptor(interceptor))
	}

	for _, interceptor := range streamServerInterceptors {
		opts = append(opts, grpc.StreamInterceptor(interceptor))
	}
	grpcServer := grpc.NewServer(opts...)

	for _, rf := range registrationFunctions {
		// i.e. protobuf.RegisterSomeExampleServiceServer(grpcServer, &exampleServiceServerImpl{})
		rf(grpcServer)
	}

	return &Server{
		logger:     logging.EnsureLogger(logger).WithName(serviceName),
		config:     cfg,
		grpcServer: grpcServer,
	}, nil
}

// Shutdown shuts down the server.
func (s *Server) Shutdown() {
	s.grpcServer.Stop()
}

// Serve serves HTTP traffic.
func (s *Server) Serve() {
	s.logger.Info("serve invoked, setting up listener")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.logger.WithValue("portx", s.config.Port).Info("listener established, serving")
	if err = s.grpcServer.Serve(lis); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop.
			os.Exit(0)
		}
	}
}
