package grpcconverters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
)

func ConvertGRPCQueryFilterToQueryFilter(qf *messages.QueryFilter) *filtering.QueryFilter {
	if qf == nil {
		return filtering.DefaultQueryFilter()
	}

	filter := &filtering.QueryFilter{
		Page:            pointer.To[uint16](1),
		Limit:           pointer.To[uint8](50),
		IncludeArchived: qf.IncludeArchived,
		// Query:           qf.Query,
	}

	return filter
}
