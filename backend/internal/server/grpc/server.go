package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/proto/service"
	"github.com/dinnerdonebetter/backend/internal/proto/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	"google.golang.org/grpc"
)

var _ service.DinnerDoneBetterServer = (*Server)(nil)

type Server struct {
	service.UnimplementedDinnerDoneBetterServer

	config *Config
	db     database.DataManager
}

func NewServer(config *Config, db database.DataManager) (*Server, error) {
	if config == nil {
		return nil, fmt.Errorf("nil config provided")
	}

	s := &Server{
		db:     db,
		config: config,
	}

	return s, nil
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %v", s.config.Port, err)
	}

	grpcServer := grpc.NewServer()
	service.RegisterDinnerDoneBetterServer(grpcServer, s)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

func (s *Server) GetValidIngredient(ctx context.Context, request *types.GetValidIngredientRequest) (*types.ValidIngredient, error) {
	validIngredient, err := s.db.GetValidIngredient(ctx, request.ValidIngredientID)
	if err != nil {
		return nil, err
	}

	return converters.ConvertValidIngredientToProtoValidIngredient(validIngredient), nil
}
