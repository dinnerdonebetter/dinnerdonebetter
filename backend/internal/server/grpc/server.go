package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/proto/service"
	"github.com/dinnerdonebetter/backend/internal/proto/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

var _ service.DinnerDoneBetterServer = (*Server)(nil)

type Server struct {
	service.UnimplementedDinnerDoneBetterServer

	db database.DataManager
}

func (s *Server) GetValidIngredient(ctx context.Context, request *types.GetValidIngredientRequest) (*types.ValidIngredient, error) {
	validIngredient, err := s.db.GetValidIngredient(ctx, request.ValidIngredientID)
	if err != nil {
		return nil, err
	}

	return converters.ConvertValidIngredientToProtoValidIngredient(validIngredient), nil
}
