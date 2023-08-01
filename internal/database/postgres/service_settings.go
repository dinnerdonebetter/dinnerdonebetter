package postgres

import (
	"context"
	_ "embed"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/database"
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

	// serviceSettingsTableColumns are the columns for the service_settings table.
	serviceSettingsTableColumns = []string{
		"service_settings.id",
		"service_settings.name",
		"service_settings.type",
		"service_settings.description",
		"service_settings.default_value",
		"service_settings.admins_only",
		"service_settings.enumeration",
		"service_settings.created_at",
		"service_settings.last_updated_at",
		"service_settings.archived_at",
	}
)

// scanServiceSetting takes a database Scanner (i.e. *sql.Row) and scans the result into a service setting struct.
func (q *Querier) scanServiceSetting(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ServiceSetting, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ServiceSetting{}

	var joinedEnum string
	targetVars := []any{
		&x.ID,
		&x.Name,
		&x.Type,
		&x.Description,
		&x.DefaultValue,
		&x.AdminsOnly,
		&joinedEnum,
		&x.CreatedAt,
		&x.LastUpdatedAt,
		&x.ArchivedAt,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	x.Enumeration = []string{}
	for _, y := range strings.Split(joinedEnum, serviceSettingsEnumDelimiter) {
		if y != "" {
			x.Enumeration = append(x.Enumeration, y)
		}
	}

	return x, filteredCount, totalCount, nil
}

// scanServiceSettings takes some database rows and turns them into a slice of service settings.
func (q *Querier) scanServiceSettings(ctx context.Context, rows database.ResultIterator, includeCounts bool) (serviceSettings []*types.ServiceSetting, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanServiceSetting(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		serviceSettings = append(serviceSettings, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return serviceSettings, filteredCount, totalCount, nil
}

//go:embed queries/service_settings/exists.sql
var serviceSettingExistenceQuery string

// ServiceSettingExists fetches whether a service setting exists from the database.
func (q *Querier) ServiceSettingExists(ctx context.Context, serviceSettingID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if serviceSettingID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingIDKey, serviceSettingID)
	tracing.AttachServiceSettingIDToSpan(span, serviceSettingID)

	args := []any{
		serviceSettingID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, serviceSettingExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing service setting existence check")
	}

	return result, nil
}

//go:embed queries/service_settings/get_one.sql
var getServiceSettingQuery string

// GetServiceSetting fetches a service setting from the database.
func (q *Querier) GetServiceSetting(ctx context.Context, serviceSettingID string) (*types.ServiceSetting, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if serviceSettingID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingIDKey, serviceSettingID)
	tracing.AttachServiceSettingIDToSpan(span, serviceSettingID)

	args := []any{
		serviceSettingID,
	}

	row := q.getOneRow(ctx, q.db, "service setting", getServiceSettingQuery, args)

	serviceSetting, _, _, err := q.scanServiceSetting(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning service setting")
	}

	return serviceSetting, nil
}

//go:embed queries/service_settings/search.sql
var serviceSettingSearchQuery string

// SearchForServiceSettings fetches a service setting from the database.
func (q *Querier) SearchForServiceSettings(ctx context.Context, query string) ([]*types.ServiceSetting, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachServiceSettingIDToSpan(span, query)

	args := []any{
		wrapQueryForILIKE(query),
	}

	rows, err := q.getRows(ctx, q.db, "service settings", serviceSettingSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing service settings list retrieval query")
	}

	x, _, _, err := q.scanServiceSettings(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning service settings")
	}

	return x, nil
}

// GetServiceSettings fetches a list of service settings from the database that meet a particular filter.
func (q *Querier) GetServiceSettings(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ServiceSetting], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.ServiceSetting]{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	query, args := q.buildListQuery(ctx, "service_settings", nil, nil, nil, householdOwnershipColumn, serviceSettingsTableColumns, "", false, filter)

	rows, err := q.getRows(ctx, q.db, "service settings", query, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing service settings list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanServiceSettings(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning service settings")
	}

	return x, nil
}

//go:embed queries/service_settings/create.sql
var serviceSettingCreationQuery string

// CreateServiceSetting creates a service setting in the database.
func (q *Querier) CreateServiceSetting(ctx context.Context, input *types.ServiceSettingDatabaseCreationInput) (*types.ServiceSetting, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ServiceSettingIDKey, input.ID)

	args := []any{
		input.ID,
		input.Name,
		input.Type,
		input.Description,
		input.DefaultValue,
		input.AdminsOnly,
		strings.Join(input.Enumeration, serviceSettingsEnumDelimiter),
	}

	// create the service setting.
	if err := q.performWriteQuery(ctx, q.db, "service setting creation", serviceSettingCreationQuery, args); err != nil {
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

	tracing.AttachServiceSettingIDToSpan(span, x.ID)
	logger.Info("service setting created")

	return x, nil
}

//go:embed queries/service_settings/archive.sql
var archiveServiceSettingQuery string

// ArchiveServiceSetting archives a service setting from the database by its ID.
func (q *Querier) ArchiveServiceSetting(ctx context.Context, serviceSettingID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if serviceSettingID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingIDKey, serviceSettingID)
	tracing.AttachServiceSettingIDToSpan(span, serviceSettingID)

	args := []any{
		serviceSettingID,
	}

	if err := q.performWriteQuery(ctx, q.db, "service setting archive", archiveServiceSettingQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating service setting")
	}

	logger.Info("service setting archived")

	return nil
}
