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

var (
	_ types.ServiceSettingConfigurationDataManager = (*Querier)(nil)
)

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

	result, err := q.generatedQuerier.CheckServiceSettingConfigurationExistence(ctx, q.db, serviceSettingConfigurationID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing service setting configuration existence check")
	}

	return result, nil
}

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

	result, err := q.generatedQuerier.GetServiceSettingConfigurationByID(ctx, q.db, serviceSettingConfigurationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching service setting configuration")
	}

	usableEnumeration := []string{}
	for _, x := range strings.Split(result.ServiceSettingEnumeration, serviceSettingsEnumDelimiter) {
		if strings.TrimSpace(x) != "" {
			usableEnumeration = append(usableEnumeration, x)
		}
	}

	serviceSettingConfiguration := &types.ServiceSettingConfiguration{
		CreatedAt:          result.CreatedAt,
		LastUpdatedAt:      timePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:         timePointerFromNullTime(result.ArchivedAt),
		ID:                 result.ID,
		Value:              result.Value,
		Notes:              result.Notes,
		BelongsToUser:      result.BelongsToUser,
		BelongsToHousehold: result.BelongsToHousehold,
		ServiceSetting: types.ServiceSetting{
			CreatedAt:     result.ServiceSettingCreatedAt,
			DefaultValue:  stringPointerFromNullString(result.ServiceSettingDefaultValue),
			LastUpdatedAt: timePointerFromNullTime(result.ServiceSettingLastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ServiceSettingArchivedAt),
			ID:            result.ServiceSettingID,
			Name:          result.ServiceSettingName,
			Type:          string(result.ServiceSettingType),
			Description:   result.ServiceSettingDescription,
			Enumeration:   usableEnumeration,
			AdminsOnly:    result.ServiceSettingAdminsOnly,
		},
	}

	return serviceSettingConfiguration, nil
}

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

	result, err := q.generatedQuerier.GetServiceSettingConfigurationForUserBySettingName(ctx, q.db, &generated.GetServiceSettingConfigurationForUserBySettingNameParams{
		Name:          settingName,
		BelongsToUser: userID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching service setting configuration")
	}

	usableEnumeration := []string{}
	for _, x := range strings.Split(result.ServiceSettingEnumeration, serviceSettingsEnumDelimiter) {
		if strings.TrimSpace(x) != "" {
			usableEnumeration = append(usableEnumeration, x)
		}
	}

	serviceSettingConfiguration := &types.ServiceSettingConfiguration{
		CreatedAt:          result.CreatedAt,
		LastUpdatedAt:      timePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:         timePointerFromNullTime(result.ArchivedAt),
		ID:                 result.ID,
		Value:              result.Value,
		Notes:              result.Notes,
		BelongsToUser:      result.BelongsToUser,
		BelongsToHousehold: result.BelongsToHousehold,
		ServiceSetting: types.ServiceSetting{
			CreatedAt:     result.ServiceSettingCreatedAt,
			DefaultValue:  stringPointerFromNullString(result.ServiceSettingDefaultValue),
			LastUpdatedAt: timePointerFromNullTime(result.ServiceSettingLastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ServiceSettingArchivedAt),
			ID:            result.ServiceSettingID,
			Name:          result.ServiceSettingName,
			Type:          string(result.ServiceSettingType),
			Description:   result.ServiceSettingDescription,
			Enumeration:   usableEnumeration,
			AdminsOnly:    result.ServiceSettingAdminsOnly,
		},
	}

	return serviceSettingConfiguration, nil
}

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

	result, err := q.generatedQuerier.GetServiceSettingConfigurationForHouseholdBySettingName(ctx, q.db, &generated.GetServiceSettingConfigurationForHouseholdBySettingNameParams{
		Name:               settingName,
		BelongsToHousehold: householdID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching service setting configuration")
	}

	usableEnumeration := []string{}
	for _, x := range strings.Split(result.ServiceSettingEnumeration, serviceSettingsEnumDelimiter) {
		if strings.TrimSpace(x) != "" {
			usableEnumeration = append(usableEnumeration, x)
		}
	}

	serviceSettingConfiguration := &types.ServiceSettingConfiguration{
		CreatedAt:          result.CreatedAt,
		LastUpdatedAt:      timePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:         timePointerFromNullTime(result.ArchivedAt),
		ID:                 result.ID,
		Value:              result.Value,
		Notes:              result.Notes,
		BelongsToUser:      result.BelongsToUser,
		BelongsToHousehold: result.BelongsToHousehold,
		ServiceSetting: types.ServiceSetting{
			CreatedAt:     result.ServiceSettingCreatedAt,
			DefaultValue:  stringPointerFromNullString(result.ServiceSettingDefaultValue),
			LastUpdatedAt: timePointerFromNullTime(result.ServiceSettingLastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ServiceSettingArchivedAt),
			ID:            result.ServiceSettingID,
			Name:          result.ServiceSettingName,
			Type:          string(result.ServiceSettingType),
			Description:   result.ServiceSettingDescription,
			Enumeration:   usableEnumeration,
			AdminsOnly:    result.ServiceSettingAdminsOnly,
		},
	}

	return serviceSettingConfiguration, nil
}

// GetServiceSettingConfigurationsForUser fetches a list of service setting configurations from the database that meet a particular filter.
func (q *Querier) GetServiceSettingConfigurationsForUser(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachUserIDToSpan(span, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	filter.AttachToLogger(logger)

	x := &types.QueryFilteredResult[types.ServiceSettingConfiguration]{
		Pagination: filter.ToPagination(),
	}

	// TODO: properly apply query filter to this
	results, err := q.generatedQuerier.GetServiceSettingConfigurationsForUser(ctx, q.db, userID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing service setting configurations list retrieval query")
	}

	for _, result := range results {
		usableEnumeration := []string{}
		for _, x := range strings.Split(result.ServiceSettingEnumeration, serviceSettingsEnumDelimiter) {
			if strings.TrimSpace(x) != "" {
				usableEnumeration = append(usableEnumeration, x)
			}
		}

		serviceSettingConfiguration := &types.ServiceSettingConfiguration{
			CreatedAt:          result.CreatedAt,
			LastUpdatedAt:      timePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:         timePointerFromNullTime(result.ArchivedAt),
			ID:                 result.ID,
			Value:              result.Value,
			Notes:              result.Notes,
			BelongsToUser:      result.BelongsToUser,
			BelongsToHousehold: result.BelongsToHousehold,
			ServiceSetting: types.ServiceSetting{
				CreatedAt:     result.ServiceSettingCreatedAt,
				DefaultValue:  stringPointerFromNullString(result.ServiceSettingDefaultValue),
				LastUpdatedAt: timePointerFromNullTime(result.ServiceSettingLastUpdatedAt),
				ArchivedAt:    timePointerFromNullTime(result.ServiceSettingArchivedAt),
				ID:            result.ServiceSettingID,
				Name:          result.ServiceSettingName,
				Type:          string(result.ServiceSettingType),
				Description:   result.ServiceSettingDescription,
				Enumeration:   usableEnumeration,
				AdminsOnly:    result.ServiceSettingAdminsOnly,
			},
		}

		x.Data = append(x.Data, serviceSettingConfiguration)
	}

	return x, nil
}

// GetServiceSettingConfigurationsForHousehold fetches a list of service setting configurations from the database that meet a particular filter.
func (q *Querier) GetServiceSettingConfigurationsForHousehold(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachHouseholdIDToSpan(span, householdID)
	logger = logger.WithValue(keys.HouseholdIDKey, householdID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	filter.AttachToLogger(logger)

	x := &types.QueryFilteredResult[types.ServiceSettingConfiguration]{
		Pagination: filter.ToPagination(),
	}

	// TODO: properly apply query filter to this
	results, err := q.generatedQuerier.GetServiceSettingConfigurationsForHousehold(ctx, q.db, householdID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing service setting configurations list retrieval query")
	}

	for _, result := range results {
		usableEnumeration := []string{}
		for _, x := range strings.Split(result.ServiceSettingEnumeration, serviceSettingsEnumDelimiter) {
			if strings.TrimSpace(x) != "" {
				usableEnumeration = append(usableEnumeration, x)
			}
		}

		serviceSettingConfiguration := &types.ServiceSettingConfiguration{
			CreatedAt:          result.CreatedAt,
			LastUpdatedAt:      timePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:         timePointerFromNullTime(result.ArchivedAt),
			ID:                 result.ID,
			Value:              result.Value,
			Notes:              result.Notes,
			BelongsToUser:      result.BelongsToUser,
			BelongsToHousehold: result.BelongsToHousehold,
			ServiceSetting: types.ServiceSetting{
				CreatedAt:     result.ServiceSettingCreatedAt,
				DefaultValue:  stringPointerFromNullString(result.ServiceSettingDefaultValue),
				LastUpdatedAt: timePointerFromNullTime(result.ServiceSettingLastUpdatedAt),
				ArchivedAt:    timePointerFromNullTime(result.ServiceSettingArchivedAt),
				ID:            result.ServiceSettingID,
				Name:          result.ServiceSettingName,
				Type:          string(result.ServiceSettingType),
				Description:   result.ServiceSettingDescription,
				Enumeration:   usableEnumeration,
				AdminsOnly:    result.ServiceSettingAdminsOnly,
			},
		}

		x.Data = append(x.Data, serviceSettingConfiguration)
	}

	return x, nil
}

// CreateServiceSettingConfiguration creates a service setting configuration in the database.
func (q *Querier) CreateServiceSettingConfiguration(ctx context.Context, input *types.ServiceSettingConfigurationDatabaseCreationInput) (*types.ServiceSettingConfiguration, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachServiceSettingConfigurationIDToSpan(span, input.ID)
	logger := q.logger.WithValue(keys.ServiceSettingConfigurationIDKey, input.ID)

	// create the service setting configuration.
	if err := q.generatedQuerier.CreateServiceSettingConfiguration(ctx, q.db, &generated.CreateServiceSettingConfigurationParams{
		ID:                 input.ID,
		Value:              input.Value,
		Notes:              input.Notes,
		ServiceSettingID:   input.ServiceSettingID,
		BelongsToUser:      input.BelongsToUser,
		BelongsToHousehold: input.BelongsToHousehold,
	}); err != nil {
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

	logger.Info("service setting configuration created")

	return x, nil
}

// UpdateServiceSettingConfiguration updates a particular service setting configuration.
func (q *Querier) UpdateServiceSettingConfiguration(ctx context.Context, updated *types.ServiceSettingConfiguration) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ServiceSettingConfigurationIDKey, updated.ID)
	tracing.AttachServiceSettingConfigurationIDToSpan(span, updated.ID)

	if err := q.generatedQuerier.UpdateServiceSettingConfiguration(ctx, q.db, &generated.UpdateServiceSettingConfigurationParams{
		Value:              updated.Value,
		Notes:              updated.Notes,
		ServiceSettingID:   updated.ServiceSetting.ID,
		BelongsToUser:      updated.BelongsToUser,
		BelongsToHousehold: updated.BelongsToHousehold,
		ID:                 updated.ID,
	}); err != nil {
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
