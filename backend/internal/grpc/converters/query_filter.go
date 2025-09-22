package grpcconverters

import (
	grpcfiltering "github.com/dinnerdonebetter/backend/internal/grpc/generated/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
)

func ConvertGRPCQueryFilterToQueryFilter(qf *grpcfiltering.QueryFilter) *filtering.QueryFilter {
	if qf == nil {
		return filtering.DefaultQueryFilter()
	}

	// TODO: better sourcing for the Page and PageSize fields
	filter := &filtering.QueryFilter{
		Page:            pointer.To[uint16](1),
		PageSize:        pointer.To[uint8](50),
		IncludeArchived: qf.IncludeArchived,
		NextCursor:      qf.NextCursor,
		Query:           "",
	}

	return filter
}
