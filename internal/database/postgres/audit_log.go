package postgres

import (
	"context"
	"encoding/json"

	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.AuditLogEntryDataManager = (*Querier)(nil)
)

// GetAuditLogEntry fetches an audit log entry from the database.
func (q *Querier) GetAuditLogEntry(ctx context.Context, auditLogEntryID, householdID string) (*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if auditLogEntryID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AuditLogEntryIDKey, auditLogEntryID)
	tracing.AttachToSpan(span, keys.AuditLogEntryIDKey, auditLogEntryID)

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	result, err := q.generatedQuerier.GetAuditLogEntry(ctx, q.db, auditLogEntryID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entry")
	}

	auditLogEntry := &types.AuditLogEntry{
		CreatedAt:          result.CreatedAt,
		BelongsToHousehold: stringPointerFromNullString(result.BelongsToHousehold),
		ID:                 result.ID,
		ResourceType:       result.ResourceType,
		RelevantID:         result.RelevantID,
		EventType:          types.AuditLogEntryEventType(result.EventType),
		BelongsToUser:      result.BelongsToUser,
	}

	if err = json.Unmarshal(result.Changes, &auditLogEntry.Changes); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
	}

	return auditLogEntry, nil
}

// GetAuditLogEntriesForUser fetches a list of audit log entries from the database that meet a particular filter.
func (q *Querier) GetAuditLogEntriesForUser(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	tracing.AttachQueryFilterToSpan(span, filter)
	x := &types.QueryFilteredResult[types.AuditLogEntry]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetAuditLogEntriesForUser(ctx, q.db, &generated.GetAuditLogEntriesForUserParams{
		BelongsToUser: userID,
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	for _, result := range results {
		auditLogEntry := &types.AuditLogEntry{
			CreatedAt:          result.CreatedAt,
			BelongsToHousehold: stringPointerFromNullString(result.BelongsToHousehold),
			ID:                 result.ID,
			ResourceType:       result.ResourceType,
			RelevantID:         result.RelevantID,
			EventType:          types.AuditLogEntryEventType(result.EventType),
			BelongsToUser:      result.BelongsToUser,
		}

		if err = json.Unmarshal(result.Changes, &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		x.Data = append(x.Data, auditLogEntry)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetAuditLogEntriesForUserAndResourceType fetches a list of audit log entries from the database that meet a particular filter.
func (q *Querier) GetAuditLogEntriesForUserAndResourceType(ctx context.Context, userID, resourceType string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if resourceType == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue("audit_log_resource_type", resourceType)
	tracing.AttachToSpan(span, "audit_log_resource_type", resourceType)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	tracing.AttachQueryFilterToSpan(span, filter)
	x := &types.QueryFilteredResult[types.AuditLogEntry]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetAuditLogEntriesForUserAndResourceType(ctx, q.db, &generated.GetAuditLogEntriesForUserAndResourceTypeParams{
		BelongsToUser: userID,
		ResourceType:  resourceType,
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	for _, result := range results {
		auditLogEntry := &types.AuditLogEntry{
			CreatedAt:          result.CreatedAt,
			BelongsToHousehold: stringPointerFromNullString(result.BelongsToHousehold),
			ID:                 result.ID,
			ResourceType:       result.ResourceType,
			RelevantID:         result.RelevantID,
			EventType:          types.AuditLogEntryEventType(result.EventType),
			BelongsToUser:      result.BelongsToUser,
		}

		if err = json.Unmarshal(result.Changes, &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		x.Data = append(x.Data, auditLogEntry)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetAuditLogEntriesForHousehold fetches a list of audit log entries from the database that meet a particular filter.
func (q *Querier) GetAuditLogEntriesForHousehold(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	tracing.AttachQueryFilterToSpan(span, filter)
	x := &types.QueryFilteredResult[types.AuditLogEntry]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetAuditLogEntriesForHousehold(ctx, q.db, &generated.GetAuditLogEntriesForHouseholdParams{
		BelongsToHousehold: nullStringFromString(householdID),
		CreatedBefore:      nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:       nullTimeFromTimePointer(filter.CreatedAfter),
		QueryOffset:        nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:         nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	for _, result := range results {
		auditLogEntry := &types.AuditLogEntry{
			CreatedAt:          result.CreatedAt,
			BelongsToHousehold: stringPointerFromNullString(result.BelongsToHousehold),
			ID:                 result.ID,
			ResourceType:       result.ResourceType,
			RelevantID:         result.RelevantID,
			EventType:          types.AuditLogEntryEventType(result.EventType),
			BelongsToUser:      result.BelongsToUser,
		}

		if err = json.Unmarshal(result.Changes, &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		x.Data = append(x.Data, auditLogEntry)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetAuditLogEntriesForHouseholdAndResourceType fetches a list of audit log entries from the database that meet a particular filter.
func (q *Querier) GetAuditLogEntriesForHouseholdAndResourceType(ctx context.Context, householdID, resourceType string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	if resourceType == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue("audit_log_resource_type", resourceType)
	tracing.AttachToSpan(span, "audit_log_resource_type", resourceType)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}

	tracing.AttachQueryFilterToSpan(span, filter)
	x := &types.QueryFilteredResult[types.AuditLogEntry]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetAuditLogEntriesForHouseholdAndResourceType(ctx, q.db, &generated.GetAuditLogEntriesForHouseholdAndResourceTypeParams{
		BelongsToHousehold: nullStringFromString(householdID),
		ResourceType:       resourceType,
		CreatedBefore:      nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:       nullTimeFromTimePointer(filter.CreatedAfter),
		QueryOffset:        nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:         nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	for _, result := range results {
		auditLogEntry := &types.AuditLogEntry{
			CreatedAt:          result.CreatedAt,
			BelongsToHousehold: stringPointerFromNullString(result.BelongsToHousehold),
			ID:                 result.ID,
			ResourceType:       result.ResourceType,
			RelevantID:         result.RelevantID,
			EventType:          types.AuditLogEntryEventType(result.EventType),
			BelongsToUser:      result.BelongsToUser,
		}

		if err = json.Unmarshal(result.Changes, &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		x.Data = append(x.Data, auditLogEntry)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateAuditLogEntry creates an audit log entry in a database.
func (q *Querier) CreateAuditLogEntry(ctx context.Context, input *types.AuditLogEntryDatabaseCreationInput) (*types.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, input.BelongsToHousehold)
	logger := q.logger.WithValue(keys.HouseholdIDKey, input.BelongsToHousehold)

	logger.Debug("CreateAuditLogEntry invoked")

	marshaledChanges, err := json.Marshal(input.Changes)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "serializing audit log change list")
	}

	if err = q.generatedQuerier.CreateAuditLogEntry(ctx, q.db, &generated.CreateAuditLogEntryParams{
		ID:                 input.ID,
		ResourceType:       input.ResourceType,
		RelevantID:         input.RelevantID,
		EventType:          generated.AuditLogEventType(input.EventType),
		Changes:            marshaledChanges,
		BelongsToUser:      input.BelongsToUser,
		BelongsToHousehold: nullStringFromStringPointer(input.BelongsToHousehold),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing webhook creation query")
	}

	x := &types.AuditLogEntry{
		ID:                 input.ID,
		Changes:            input.Changes,
		BelongsToHousehold: input.BelongsToHousehold,
		CreatedAt:          q.currentTime(),
		ResourceType:       input.ResourceType,
		RelevantID:         input.RelevantID,
		EventType:          input.EventType,
		BelongsToUser:      input.BelongsToUser,
	}

	tracing.AttachToSpan(span, keys.AuditLogEntryIDKey, x.ID)

	return x, nil
}
