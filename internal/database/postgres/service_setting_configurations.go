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

var (
	_ types.ServiceSettingConfigurationDataManager = (*Querier)(nil)

	// serviceSettingConfigurationsTableColumns are the columns for the service_setting_configurations table.
	serviceSettingConfigurationsTableColumns = []string{
		"service_setting_configurations.id",
		"service_setting_configurations.value",
		"service_setting_configurations.notes",
		"service_settings.id",
		"service_settings.name",
		"service_settings.type",
		"service_settings.description",
		"service_settings.default_value",
		"service_settings.enumeration",
		"service_settings.admins_only",
		"service_settings.created_at",
		"service_settings.last_updated_at",
		"service_settings.archived_at",
		"service_setting_configurations.belongs_to_user",
		"service_setting_configurations.belongs_to_household",
		"service_setting_configurations.created_at",
		"service_setting_configurations.last_updated_at",
		"service_setting_configurations.archived_at",
	}
)

// scanServiceSettingConfiguration takes a database Scanner (i.e. *sql.Row) and scans the result into a service setting configuration struct.
func (q *Querier) scanServiceSettingConfiguration(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ServiceSettingConfiguration, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	x = &types.ServiceSettingConfiguration{}

	var joinedEnums string
	targetVars := []any{
		&x.ID,
		&x.Value,
		&x.Notes,
		&x.ServiceSetting.ID,
		&x.ServiceSetting.Name,
		&x.ServiceSetting.Type,
		&x.ServiceSetting.Description,
		&x.ServiceSetting.DefaultValue,
		&joinedEnums,
		&x.ServiceSetting.AdminsOnly,
		&x.ServiceSetting.CreatedAt,
		&x.ServiceSetting.LastUpdatedAt,
		&x.ServiceSetting.ArchivedAt,
		&x.BelongsToUser,
		&x.BelongsToHousehold,
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

	x.ServiceSetting.Enumeration = []string{}
	for _, y := range strings.Split(joinedEnums, serviceSettingsEnumDelimiter) {
		if y != "" {
			x.ServiceSetting.Enumeration = append(x.ServiceSetting.Enumeration, y)
		}
	}

	return x, filteredCount, totalCount, nil
}

// scanServiceSettingConfigurations takes some database rows and turns them into a slice of service setting configurations.
func (q *Querier) scanServiceSettingConfigurations(ctx context.Context, rows database.ResultIterator, includeCounts bool) (serviceSettingConfigurations []*types.ServiceSettingConfiguration, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanServiceSettingConfiguration(ctx, rows, includeCounts)
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

		serviceSettingConfigurations = append(serviceSettingConfigurations, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return serviceSettingConfigurations, filteredCount, totalCount, nil
}

//go:embed queries/service_setting_configurations/exists.sql
var serviceSettingConfigurationExistenceQuery string

// ServiceSettingConfigurationExists fetches whether a service setting configuration exists from the database.
func (q *Querier) ServiceSettingConfigurationExists(ctx context.Context, serviceSettingConfigurationID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if serviceSettingConfigurationID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)
	tracing.AttachServiceSettingConfigurationIDToSpan(span, serviceSettingConfigurationID)

	args := []any{
		serviceSettingConfigurationID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, serviceSettingConfigurationExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing service setting configuration existence check")
	}

	return result, nil
}

//go:embed queries/service_setting_configurations/get_by_id.sql
var getServiceSettingConfigurationQuery string

// GetServiceSettingConfiguration fetches a service setting configuration from the database.
func (q *Querier) GetServiceSettingConfiguration(ctx context.Context, serviceSettingConfigurationID string) (*types.ServiceSettingConfiguration, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if serviceSettingConfigurationID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachServiceSettingConfigurationIDToSpan(span, serviceSettingConfigurationID)
	logger = logger.WithValue(keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)

	args := []any{
		serviceSettingConfigurationID,
	}

	row := q.getOneRow(ctx, q.db, "service settings for user by name", getServiceSettingConfigurationQuery, args)

	serviceSettingConfiguration, _, _, err := q.scanServiceSettingConfiguration(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning service setting configuration")
	}

	return serviceSettingConfiguration, nil
}

//go:embed queries/service_setting_configurations/get_for_user_by_setting_name.sql
var getServiceSettingConfigurationForUserByNameQuery string

// GetServiceSettingConfigurationForUserByName fetches a service setting configuration from the database.
func (q *Querier) GetServiceSettingConfigurationForUserByName(ctx context.Context, userID, settingName string) (*types.ServiceSettingConfiguration, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	if settingName == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingNameKey, settingName)
	tracing.AttachServiceSettingNameToSpan(span, settingName)

	args := []any{
		settingName,
		userID,
	}

	row := q.getOneRow(ctx, q.db, "service settings for user by name", getServiceSettingConfigurationForUserByNameQuery, args)

	serviceSettingConfiguration, _, _, err := q.scanServiceSettingConfiguration(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning service setting configuration for user by name")
	}

	return serviceSettingConfiguration, nil
}

//go:embed queries/service_setting_configurations/get_for_household_by_setting_name.sql
var getServiceSettingConfigurationForHouseholdByNameQuery string

// GetServiceSettingConfigurationForHouseholdByName fetches a service setting configuration from the database.
func (q *Querier) GetServiceSettingConfigurationForHouseholdByName(ctx context.Context, householdID, settingName string) (*types.ServiceSettingConfiguration, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachHouseholdIDToSpan(span, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	if settingName == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingNameKey, settingName)
	tracing.AttachServiceSettingNameToSpan(span, settingName)

	args := []any{
		settingName,
		householdID,
	}

	row := q.getOneRow(ctx, q.db, "service settings for household by name", getServiceSettingConfigurationForHouseholdByNameQuery, args)

	serviceSettingConfiguration, _, _, err := q.scanServiceSettingConfiguration(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning service setting configuration for household by name")
	}

	return serviceSettingConfiguration, nil
}

//go:embed queries/service_setting_configurations/get_settings_for_user.sql
var getServiceSettingConfigurationForUserQuery string

// GetServiceSettingConfigurationsForUser fetches a list of service setting configurations from the database that meet a particular filter.
func (q *Querier) GetServiceSettingConfigurationsForUser(ctx context.Context, userID string) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	args := []any{
		userID,
	}

	rows, err := q.getRows(ctx, q.db, "service setting configurations for user", getServiceSettingConfigurationForUserQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing service setting configurations list retrieval query")
	}

	x := &types.QueryFilteredResult[types.ServiceSettingConfiguration]{
		Pagination: types.DefaultQueryFilter().ToPagination(),
	}
	if x.Data, _, _, err = q.scanServiceSettingConfigurations(ctx, rows, false); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning service setting configurations")
	}

	return x, nil
}

//go:embed queries/service_setting_configurations/get_settings_for_household.sql
var getServiceSettingConfigurationForHouseholdQuery string

// GetServiceSettingConfigurationsForHousehold fetches a list of service setting configurations from the database that meet a particular filter.
func (q *Querier) GetServiceSettingConfigurationsForHousehold(ctx context.Context, householdID string) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachHouseholdIDToSpan(span, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	args := []any{
		householdID,
	}

	rows, err := q.getRows(ctx, q.db, "service setting configurations for household", getServiceSettingConfigurationForHouseholdQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing service setting configurations list retrieval query")
	}

	x := &types.QueryFilteredResult[types.ServiceSettingConfiguration]{
		Pagination: types.DefaultQueryFilter().ToPagination(),
	}
	if x.Data, _, _, err = q.scanServiceSettingConfigurations(ctx, rows, false); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning service setting configurations")
	}

	return x, nil
}

//go:embed queries/service_setting_configurations/create.sql
var serviceSettingConfigurationCreationQuery string

// CreateServiceSettingConfiguration creates a service setting configuration in the database.
func (q *Querier) CreateServiceSettingConfiguration(ctx context.Context, input *types.ServiceSettingConfigurationDatabaseCreationInput) (*types.ServiceSettingConfiguration, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ServiceSettingConfigurationIDKey, input.ID)

	args := []any{
		input.ID,
		input.Value,
		input.Notes,
		input.ServiceSettingID,
		input.BelongsToUser,
		input.BelongsToHousehold,
	}

	// create the service setting configuration.
	if err := q.performWriteQuery(ctx, q.db, "service setting configuration creation", serviceSettingConfigurationCreationQuery, args); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing service setting configuration creation query")
	}

	x := &types.ServiceSettingConfiguration{
		ID:                 input.ID,
		Value:              input.Value,
		Notes:              input.Notes,
		ServiceSetting:     types.ServiceSetting{ID: input.ServiceSettingID},
		BelongsToUser:      input.BelongsToUser,
		BelongsToHousehold: input.BelongsToHousehold,
		CreatedAt:          q.currentTime(),
	}

	tracing.AttachServiceSettingConfigurationIDToSpan(span, x.ID)
	logger.Info("service setting configuration created")

	return x, nil
}

//go:embed queries/service_setting_configurations/update.sql
var updateServiceSettingConfigurationQuery string

// UpdateServiceSettingConfiguration updates a particular service setting configuration.
func (q *Querier) UpdateServiceSettingConfiguration(ctx context.Context, updated *types.ServiceSettingConfiguration) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ServiceSettingConfigurationIDKey, updated.ID)
	tracing.AttachServiceSettingConfigurationIDToSpan(span, updated.ID)

	args := []any{
		updated.Value,
		updated.Notes,
		updated.ServiceSetting.ID,
		updated.BelongsToUser,
		updated.BelongsToHousehold,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "service setting configuration update", updateServiceSettingConfigurationQuery, args); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating service setting configuration")
	}

	logger.Info("service setting configuration updated")

	return nil
}

// ArchiveServiceSettingConfiguration archives a service setting configuration from the database by its ID.
func (q *Querier) ArchiveServiceSettingConfiguration(ctx context.Context, serviceSettingConfigurationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if serviceSettingConfigurationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)
	tracing.AttachServiceSettingConfigurationIDToSpan(span, serviceSettingConfigurationID)

	if err := q.generatedQuerier.ArchiveServiceSettingConfiguration(ctx, q.db, serviceSettingConfigurationID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving service setting configuration")
	}

	logger.Info("service setting configuration archived")

	return nil
}
