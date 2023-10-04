package postgres

import (
	"context"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceSettingsEnumDelimiter = "|"
)

var (
	_ types.ServiceSettingDataManager = (*Querier)(nil)
)

// ServiceSettingExists fetches whether a service setting exists from the database.
func (q *Querier) ServiceSettingExists(ctx context.Context, serviceSettingID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if serviceSettingID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingIDKey, serviceSettingID)
	tracing.AttachToSpan(span, keys.ServiceSettingIDKey, serviceSettingID)

	result, err := q.generatedQuerier.CheckServiceSettingExistence(ctx, q.db, serviceSettingID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing service setting existence check")
	}

	return result, nil
}

// GetServiceSetting fetches a service setting from the database.
func (q *Querier) GetServiceSetting(ctx context.Context, serviceSettingID string) (*types.ServiceSetting, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if serviceSettingID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingIDKey, serviceSettingID)
	tracing.AttachToSpan(span, keys.ServiceSettingIDKey, serviceSettingID)

	result, err := q.generatedQuerier.GetServiceSetting(ctx, q.db, serviceSettingID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing service setting fetch")
	}

	usableEnumeration := []string{}
	for _, x := range strings.Split(result.Enumeration, serviceSettingsEnumDelimiter) {
		if strings.TrimSpace(x) != "" {
			usableEnumeration = append(usableEnumeration, x)
		}
	}

	serviceSetting := &types.ServiceSetting{
		CreatedAt:     result.CreatedAt,
		DefaultValue:  stringPointerFromNullString(result.DefaultValue),
		LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    timePointerFromNullTime(result.ArchivedAt),
		ID:            result.ID,
		Name:          result.Name,
		Type:          string(result.Type),
		Description:   result.Description,
		Enumeration:   usableEnumeration,
		AdminsOnly:    result.AdminsOnly,
	}

	return serviceSetting, nil
}

// SearchForServiceSettings fetches a service setting from the database.
func (q *Querier) SearchForServiceSettings(ctx context.Context, query string) ([]*types.ServiceSetting, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.ServiceSettingIDKey, query)

	results, err := q.generatedQuerier.SearchForServiceSettings(ctx, q.db, query)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing service settings list retrieval query")
	}

	x := []*types.ServiceSetting{}
	for _, result := range results {
		rawEnumeration := strings.Split(result.Enumeration, serviceSettingsEnumDelimiter)
		usableEnumeration := []string{}
		for _, y := range rawEnumeration {
			if strings.TrimSpace(y) != "" {
				usableEnumeration = append(usableEnumeration, y)
			}
		}

		serviceSetting := &types.ServiceSetting{
			CreatedAt:     result.CreatedAt,
			DefaultValue:  stringPointerFromNullString(result.DefaultValue),
			LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ArchivedAt),
			ID:            result.ID,
			Name:          result.Name,
			Type:          string(result.Type),
			Description:   result.Description,
			Enumeration:   usableEnumeration,
			AdminsOnly:    result.AdminsOnly,
		}

		x = append(x, serviceSetting)
	}

	return x, nil
}

// GetServiceSettings fetches a list of service settings from the database that meet a particular filter.
func (q *Querier) GetServiceSettings(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ServiceSetting], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ServiceSetting]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetServiceSettings(ctx, q.db, &generated.GetServiceSettingsParams{
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing service settings list retrieval query")
	}

	for _, result := range results {
		rawEnumeration := strings.Split(result.Enumeration, serviceSettingsEnumDelimiter)
		usableEnumeration := []string{}
		for _, y := range rawEnumeration {
			if strings.TrimSpace(y) != "" {
				usableEnumeration = append(usableEnumeration, y)
			}
		}

		x.Data = append(x.Data, &types.ServiceSetting{
			CreatedAt:     result.CreatedAt,
			DefaultValue:  stringPointerFromNullString(result.DefaultValue),
			LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ArchivedAt),
			ID:            result.ID,
			Name:          result.Name,
			Type:          string(result.Type),
			Description:   result.Description,
			Enumeration:   usableEnumeration,
			AdminsOnly:    result.AdminsOnly,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateServiceSetting creates a service setting in the database.
func (q *Querier) CreateServiceSetting(ctx context.Context, input *types.ServiceSettingDatabaseCreationInput) (*types.ServiceSetting, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ServiceSettingIDKey, input.ID)
	logger := q.logger.WithValue(keys.ServiceSettingIDKey, input.ID)

	// create the service setting.
	if err := q.generatedQuerier.CreateServiceSetting(ctx, q.db, &generated.CreateServiceSettingParams{
		ID:           input.ID,
		Name:         input.Name,
		Type:         generated.SettingType(input.Type),
		Description:  input.Description,
		Enumeration:  strings.Join(input.Enumeration, serviceSettingsEnumDelimiter),
		DefaultValue: nullStringFromStringPointer(input.DefaultValue),
		AdminsOnly:   input.AdminsOnly,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing service setting creation query")
	}

	x := &types.ServiceSetting{
		ID:           input.ID,
		Name:         input.Name,
		Type:         input.Type,
		Description:  input.Description,
		DefaultValue: input.DefaultValue,
		AdminsOnly:   input.AdminsOnly,
		Enumeration:  input.Enumeration,
		CreatedAt:    q.currentTime(),
	}

	logger.Info("service setting created")

	return x, nil
}

// ArchiveServiceSetting archives a service setting from the database by its ID.
func (q *Querier) ArchiveServiceSetting(ctx context.Context, serviceSettingID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if serviceSettingID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingIDKey, serviceSettingID)
	tracing.AttachToSpan(span, keys.ServiceSettingIDKey, serviceSettingID)

	if _, err := q.generatedQuerier.ArchiveServiceSetting(ctx, q.db, serviceSettingID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating service setting")
	}

	logger.Info("service setting archived")

	return nil
}
