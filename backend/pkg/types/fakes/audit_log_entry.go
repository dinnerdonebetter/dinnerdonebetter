package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// BuildFakeAuditLogEntry builds a faked valid instrument.
func BuildFakeAuditLogEntry() *types.AuditLogEntry {
	return &types.AuditLogEntry{
		ID:                 BuildFakeID(),
		CreatedAt:          BuildFakeTime(),
		Changes:            nil,
		BelongsToHousehold: nil,
		ResourceType:       "example",
		RelevantID:         BuildFakeID(),
		EventType:          types.AuditLogEventTypeOther,
		BelongsToUser:      BuildFakeID(),
	}
}

// BuildFakeAuditLogEntriesList builds a faked AuditLogEntryList.
func BuildFakeAuditLogEntriesList() *filtering.QueryFilteredResult[types.AuditLogEntry] {
	var examples []*types.AuditLogEntry
	for range exampleQuantity {
		examples = append(examples, BuildFakeAuditLogEntry())
	}

	return &filtering.QueryFilteredResult[types.AuditLogEntry]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}
