package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpcimpl/converters"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateValidVessel(ctx context.Context, input *messages.ValidVesselCreationRequestInput) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidVessel(ctx context.Context, request *messages.GetValidVesselRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidVessels(ctx context.Context, request *messages.GetValidVesselsRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRandomValidVessel(ctx context.Context, _ *emptypb.Empty) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	ingredient, err := s.dataManager.GetRandomValidVessel(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting random valid ingredient")
	}

	output := &messages.ValidVessel{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(ingredient.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(ingredient.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(ingredient.ArchivedAt),
		CapacityUnit: &messages.ValidMeasurementUnit{
			CreatedAt:     converters.ConvertTimeToPBTimestamp(ingredient.CapacityUnit.CreatedAt),
			LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(ingredient.CapacityUnit.LastUpdatedAt),
			ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(ingredient.CapacityUnit.ArchivedAt),
			Name:          ingredient.CapacityUnit.Name,
			IconPath:      ingredient.CapacityUnit.IconPath,
			ID:            ingredient.CapacityUnit.ID,
			Description:   ingredient.CapacityUnit.Description,
			PluralName:    ingredient.CapacityUnit.PluralName,
			Slug:          ingredient.CapacityUnit.Slug,
			Volumetric:    ingredient.CapacityUnit.Volumetric,
			Universal:     ingredient.CapacityUnit.Universal,
			Metric:        ingredient.CapacityUnit.Metric,
			Imperial:      ingredient.CapacityUnit.Imperial,
		},
		IconPath:                       ingredient.IconPath,
		PluralName:                     ingredient.PluralName,
		Description:                    ingredient.Description,
		Name:                           ingredient.Name,
		Slug:                           ingredient.Slug,
		Shape:                          ingredient.Shape,
		ID:                             ingredient.ID,
		WidthInMillimeters:             ingredient.WidthInMillimeters,
		LengthInMillimeters:            ingredient.LengthInMillimeters,
		HeightInMillimeters:            ingredient.HeightInMillimeters,
		Capacity:                       ingredient.Capacity,
		IncludeInGeneratedInstructions: ingredient.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          ingredient.DisplayInSummaryLists,
		UsableForStorage:               ingredient.UsableForStorage,
	}

	return output, nil
}

func (s *Server) SearchForValidVessels(ctx context.Context, request *messages.SearchForValidVesselsRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidVessel(ctx context.Context, request *messages.UpdateValidVesselRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidVessel(ctx context.Context, request *messages.ArchiveValidVesselRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
