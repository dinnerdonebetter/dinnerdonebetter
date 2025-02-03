package grpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"

	"google.golang.org/grpc"
)

type (
	Config struct {
		Port uint16 `env:"PORT" json:"port"`
	}

	Server struct {
		config     *Config
		grpcServer *grpc.Server
	}
)

func NewGRPCServer(cfg *Config, registrationFunctions ...func(server *grpc.Server)) (*Server, error) {
	if cfg == nil {
		return nil, internalerrors.NilConfigError("grpc server")
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	for _, rf := range registrationFunctions {
		// i.e. protobuf.RegisterSomeExampleServiceServer(grpcServer, &exampleServiceServerImpl{})
		rf(grpcServer)
	}

	return &Server{
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
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err = s.grpcServer.Serve(lis); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop.
			os.Exit(0)
		}
	}
}
