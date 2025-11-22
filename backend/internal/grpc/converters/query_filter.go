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

	filter := &filtering.QueryFilter{
		Page:            pointer.To[uint16](1),
		PageSize:        pointer.To[uint8](50),
		IncludeArchived: qf.IncludeArchived,
		NextCursor:      qf.NextCursor,
		Query:           "",
		//
		SortBy:        nil,
		CreatedAfter:  nil,
		CreatedBefore: nil,
		UpdatedAfter:  nil,
		UpdatedBefore: nil,
	}
	if qf.PageSize != nil {
		filter.PageSize = pointer.To(uint8(*qf.PageSize))
	}

	return filter
}

func ConvertQueryFilterToGRPCQueryFilter(qf *filtering.QueryFilter, resultPagination filtering.Pagination) *grpcfiltering.QueryFilter {
	if qf == nil {
		qf = filtering.DefaultQueryFilter()
	}

	f := &grpcfiltering.QueryFilter{
		IncludeArchived: qf.IncludeArchived,
		NextCursor:      qf.NextCursor,
	}
	if qf.PageSize != nil {
		f.PageSize = pointer.To(uint32(*qf.PageSize))
	}

	return f
}
