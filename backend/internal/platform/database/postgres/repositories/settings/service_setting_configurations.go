package settings

import (
	"context"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/domain/auditlogentries"
	types "github.com/dinnerdonebetter/backend/internal/domain/settings"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/repositories/settings/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	resourceTypeServiceSettingConfigurations = "service_setting_configurations"
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
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)
	tracing.AttachToSpan(span, keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)

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
		return nil, database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)
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
		CreatedAt:        result.CreatedAt,
		LastUpdatedAt:    database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:       database.TimePointerFromNullTime(result.ArchivedAt),
		ID:               result.ID,
		Value:            result.Value,
		Notes:            result.Notes,
		BelongsToUser:    result.BelongsToUser,
		BelongsToAccount: result.BelongsToAccount,
		ServiceSetting: types.ServiceSetting{
			CreatedAt:     result.ServiceSettingCreatedAt,
			DefaultValue:  database.StringPointerFromNullString(result.ServiceSettingDefaultValue),
			LastUpdatedAt: database.TimePointerFromNullTime(result.ServiceSettingLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ServiceSettingArchivedAt),
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
		return nil, database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	if settingName == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingNameKey, settingName)
	tracing.AttachToSpan(span, keys.ServiceSettingNameKey, settingName)

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
		CreatedAt:        result.CreatedAt,
		LastUpdatedAt:    database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:       database.TimePointerFromNullTime(result.ArchivedAt),
		ID:               result.ID,
		Value:            result.Value,
		Notes:            result.Notes,
		BelongsToUser:    result.BelongsToUser,
		BelongsToAccount: result.BelongsToAccount,
		ServiceSetting: types.ServiceSetting{
			CreatedAt:     result.ServiceSettingCreatedAt,
			DefaultValue:  database.StringPointerFromNullString(result.ServiceSettingDefaultValue),
			LastUpdatedAt: database.TimePointerFromNullTime(result.ServiceSettingLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ServiceSettingArchivedAt),
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

// GetServiceSettingConfigurationForAccountByName fetches a service setting configuration from the database.
func (q *Querier) GetServiceSettingConfigurationForAccountByName(ctx context.Context, accountID, settingName string) (*types.ServiceSettingConfiguration, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if accountID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)
	logger = logger.WithValue(keys.AccountIDKey, accountID)

	if settingName == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingNameKey, settingName)
	tracing.AttachToSpan(span, keys.ServiceSettingNameKey, settingName)

	result, err := q.generatedQuerier.GetServiceSettingConfigurationForAccountBySettingName(ctx, q.db, &generated.GetServiceSettingConfigurationForAccountBySettingNameParams{
		Name:             settingName,
		BelongsToAccount: accountID,
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
		CreatedAt:        result.CreatedAt,
		LastUpdatedAt:    database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:       database.TimePointerFromNullTime(result.ArchivedAt),
		ID:               result.ID,
		Value:            result.Value,
		Notes:            result.Notes,
		BelongsToUser:    result.BelongsToUser,
		BelongsToAccount: result.BelongsToAccount,
		ServiceSetting: types.ServiceSetting{
			CreatedAt:     result.ServiceSettingCreatedAt,
			DefaultValue:  database.StringPointerFromNullString(result.ServiceSettingDefaultValue),
			LastUpdatedAt: database.TimePointerFromNullTime(result.ServiceSettingLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ServiceSettingArchivedAt),
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
func (q *Querier) GetServiceSettingConfigurationsForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger = logger.WithValue(keys.UserIDKey, userID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	filter.AttachToLogger(logger)

	x := &filtering.QueryFilteredResult[types.ServiceSettingConfiguration]{
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
			CreatedAt:        result.CreatedAt,
			LastUpdatedAt:    database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:       database.TimePointerFromNullTime(result.ArchivedAt),
			ID:               result.ID,
			Value:            result.Value,
			Notes:            result.Notes,
			BelongsToUser:    result.BelongsToUser,
			BelongsToAccount: result.BelongsToAccount,
			ServiceSetting: types.ServiceSetting{
				CreatedAt:     result.ServiceSettingCreatedAt,
				DefaultValue:  database.StringPointerFromNullString(result.ServiceSettingDefaultValue),
				LastUpdatedAt: database.TimePointerFromNullTime(result.ServiceSettingLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ServiceSettingArchivedAt),
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

// GetServiceSettingConfigurationsForAccount fetches a list of service setting configurations from the database that meet a particular filter.
func (q *Querier) GetServiceSettingConfigurationsForAccount(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if accountID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)
	logger = logger.WithValue(keys.AccountIDKey, accountID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	tracing.AttachQueryFilterToSpan(span, filter)
	filter.AttachToLogger(logger)

	x := &filtering.QueryFilteredResult[types.ServiceSettingConfiguration]{
		Pagination: filter.ToPagination(),
	}

	// TODO: properly apply query filter to this
	results, err := q.generatedQuerier.GetServiceSettingConfigurationsForAccount(ctx, q.db, accountID)
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
			CreatedAt:        result.CreatedAt,
			LastUpdatedAt:    database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:       database.TimePointerFromNullTime(result.ArchivedAt),
			ID:               result.ID,
			Value:            result.Value,
			Notes:            result.Notes,
			BelongsToUser:    result.BelongsToUser,
			BelongsToAccount: result.BelongsToAccount,
			ServiceSetting: types.ServiceSetting{
				CreatedAt:     result.ServiceSettingCreatedAt,
				DefaultValue:  database.StringPointerFromNullString(result.ServiceSettingDefaultValue),
				LastUpdatedAt: database.TimePointerFromNullTime(result.ServiceSettingLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ServiceSettingArchivedAt),
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
		return nil, database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ServiceSettingConfigurationIDKey, input.ID)
	logger := q.logger.WithValue(keys.ServiceSettingConfigurationIDKey, input.ID)

	// begin account creation transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	// create the service setting configuration.
	if err = q.generatedQuerier.CreateServiceSettingConfiguration(ctx, q.db, &generated.CreateServiceSettingConfigurationParams{
		ID:               input.ID,
		Value:            input.Value,
		Notes:            input.Notes,
		ServiceSettingID: input.ServiceSettingID,
		BelongsToUser:    input.BelongsToUser,
		BelongsToAccount: input.BelongsToAccount,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing service setting configuration creation query")
	}

	x := &types.ServiceSettingConfiguration{
		ID:               input.ID,
		Value:            input.Value,
		Notes:            input.Notes,
		ServiceSetting:   types.ServiceSetting{ID: input.ServiceSettingID},
		BelongsToUser:    input.BelongsToUser,
		BelongsToAccount: input.BelongsToAccount,
		CreatedAt:        q.CurrentTime(),
	}

	if _, err = q.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &auditlogentries.AuditLogEntryDatabaseCreationInput{
		BelongsToAccount: &input.BelongsToAccount,
		ID:               identifiers.New(),
		ResourceType:     resourceTypeServiceSettingConfigurations,
		RelevantID:       x.ID,
		EventType:        auditlogentries.AuditLogEventTypeCreated,
		BelongsToUser:    input.BelongsToUser,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("service setting configuration created")

	return x, nil
}

// UpdateServiceSettingConfiguration updates a particular service setting configuration.
func (q *Querier) UpdateServiceSettingConfiguration(ctx context.Context, updated *types.ServiceSettingConfiguration) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ServiceSettingConfigurationIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ServiceSettingConfigurationIDKey, updated.ID)

	// begin account creation transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if _, err = q.generatedQuerier.UpdateServiceSettingConfiguration(ctx, q.db, &generated.UpdateServiceSettingConfigurationParams{
		Value:            updated.Value,
		Notes:            updated.Notes,
		ServiceSettingID: updated.ServiceSetting.ID,
		BelongsToUser:    updated.BelongsToUser,
		BelongsToAccount: updated.BelongsToAccount,
		ID:               updated.ID,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "updating service setting configuration")
	}

	if _, err = q.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &auditlogentries.AuditLogEntryDatabaseCreationInput{
		BelongsToAccount: &updated.BelongsToAccount,
		ID:               identifiers.New(),
		ResourceType:     resourceTypeServiceSettingConfigurations,
		RelevantID:       updated.ID,
		EventType:        auditlogentries.AuditLogEventTypeUpdated,
		BelongsToUser:    updated.BelongsToUser,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
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
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)
	tracing.AttachToSpan(span, keys.ServiceSettingConfigurationIDKey, serviceSettingConfigurationID)

	// begin account creation transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if _, err = q.generatedQuerier.ArchiveServiceSettingConfiguration(ctx, q.db, serviceSettingConfigurationID); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving service setting configuration")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("service setting configuration archived")

	return nil
}
