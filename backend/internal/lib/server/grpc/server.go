package grpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	serviceName = "grpc_server"
)

type (
	Config struct {
		HTTPSCertificateFile  string `env:"TLS_CERTIFICATE_FILEPATH"     json:"tlsCertificate,omitempty"`
		TLSCertificateKeyFile string `env:"TLS_CERTIFICATE_KEY_FILEPATH" json:"tlsCertificateKey,omitempty"`
		Port                  uint16 `env:"PORT"                         json:"port"`
	}

	Server struct {
		logger     logging.Logger
		config     *Config
		grpcServer *grpc.Server
	}

	// RegistrationFunc is i.e. protobuf.RegisterSomeExampleServiceServer(grpcServer, &exampleServiceServerImpl{}).
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

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(append([]grpc.UnaryServerInterceptor{LoggingInterceptor(logger)}, unaryServerInterceptors...)...),
		grpc.ChainStreamInterceptor(streamServerInterceptors...),
	}

	if cfg.TLSCertificateKeyFile != "" && cfg.HTTPSCertificateFile != "" {
		serverCert, err := tls.LoadX509KeyPair(cfg.HTTPSCertificateFile, cfg.TLSCertificateKeyFile)
		if err != nil {
			return nil, err
		}

		// Create the credentials and return it
		config := &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.NoClientCert,
		}

		opts = append(opts, grpc.Creds(credentials.NewTLS(config)))
	}

	grpcServer := grpc.NewServer(opts...)
	for _, rf := range registrationFunctions {
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

	s.logger.WithValue("port", s.config.Port).Info("listener established, serving")
	if err = s.grpcServer.Serve(lis); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			// NOTE: there is a chance that next line won't have tim  e to run,
			// as main() doesn't wait for this goroutine to stop.
			os.Exit(0)
		}
	}
}

func LoggingInterceptor(logger logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		logger.WithValue("rpc.method", info.FullMethod).Info("rpc invoked")
		return handler(ctx, req)
	}
}
