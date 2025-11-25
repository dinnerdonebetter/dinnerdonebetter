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
		Limit:           pointer.To[uint8](50),
		IncludeArchived: qf.IncludeArchived,
		Cursor:          qf.Cursor,
		SortBy:          qf.SortBy,
		CreatedAfter:    ConvertPBTimestampToTimePointer(qf.CreatedAfter),
		CreatedBefore:   ConvertPBTimestampToTimePointer(qf.CreatedBefore),
		UpdatedAfter:    ConvertPBTimestampToTimePointer(qf.UpdatedAfter),
		UpdatedBefore:   ConvertPBTimestampToTimePointer(qf.UpdatedBefore),
	}
	if qf.PageSize != nil {
		filter.Limit = pointer.To(uint8(*qf.PageSize))
	}

	return filter
}

func ConvertQueryFilterToGRPCQueryFilter(qf *filtering.QueryFilter, resultPagination filtering.Pagination) *grpcfiltering.QueryFilter {
	if qf == nil {
		qf = filtering.DefaultQueryFilter()
	}

	// Use the cursor from the result pagination if available (for next page), otherwise use the filter's cursor
	cursor := qf.Cursor
	if resultPagination.Cursor != "" {
		cursor = pointer.To(resultPagination.Cursor)
	}

	f := &grpcfiltering.QueryFilter{
		IncludeArchived: qf.IncludeArchived,
		Cursor:          cursor,
		SortBy:          qf.SortBy,
		CreatedAfter:    ConvertTimePointerToPBTimestamp(qf.CreatedAfter),
		CreatedBefore:   ConvertTimePointerToPBTimestamp(qf.CreatedBefore),
		UpdatedAfter:    ConvertTimePointerToPBTimestamp(qf.UpdatedAfter),
		UpdatedBefore:   ConvertTimePointerToPBTimestamp(qf.UpdatedBefore),
	}
	if qf.Limit != nil {
		f.PageSize = pointer.To(uint32(*qf.Limit))
	}

	return f
}

func ConvertPaginationToGRPCPagination(pagination filtering.Pagination, qf *filtering.QueryFilter) *grpcfiltering.Pagination {
	if qf == nil {
		qf = filtering.DefaultQueryFilter()
	}

	f := &grpcfiltering.Pagination{
		AppliedQueryFilter: ConvertQueryFilterToGRPCQueryFilter(qf, pagination),
		Cursor:             pagination.Cursor,
		FilteredCount:      pagination.FilteredCount,
		TotalCount:         pagination.TotalCount,
		MaxResponseSize:    uint32(pagination.MaxResponseSize),
	}

	return f
}
