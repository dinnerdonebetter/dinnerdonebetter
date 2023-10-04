package postgres

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.HouseholdDataManager = (*Querier)(nil)
)

// GetHousehold fetches a household from the database.
func (q *Querier) GetHousehold(ctx context.Context, householdID string) (*types.Household, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	results, err := q.generatedQuerier.GetHouseholdByIDWithMemberships(ctx, q.db, householdID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing households list retrieval query")
	}

	var household *types.Household
	for _, result := range results {
		if household == nil {
			household = &types.Household{
				CreatedAt:                  result.CreatedAt,
				SubscriptionPlanID:         stringPointerFromNullString(result.SubscriptionPlanID),
				LastUpdatedAt:              timePointerFromNullTime(result.LastUpdatedAt),
				ArchivedAt:                 timePointerFromNullTime(result.ArchivedAt),
				ContactPhone:               result.ContactPhone,
				BillingStatus:              result.BillingStatus,
				AddressLine1:               result.AddressLine1,
				AddressLine2:               result.AddressLine2,
				City:                       result.City,
				State:                      result.State,
				ZipCode:                    result.ZipCode,
				Country:                    result.Country,
				Latitude:                   float64PointerFromNullString(result.Latitude),
				Longitude:                  float64PointerFromNullString(result.Longitude),
				PaymentProcessorCustomerID: result.PaymentProcessorCustomerID,
				BelongsToUser:              result.BelongsToUser,
				ID:                         result.ID,
				Name:                       result.Name,
				WebhookEncryptionKey:       result.WebhookHmacSecret,
				Members:                    nil,
			}
		}

		household.Members = append(household.Members, &types.HouseholdUserMembershipWithUser{
			CreatedAt:     result.MembershipCreatedAt,
			LastUpdatedAt: timePointerFromNullTime(result.MembershipLastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.MembershipArchivedAt),
			ID:            result.MembershipID,
			BelongsToUser: &types.User{
				CreatedAt:                  result.UserCreatedAt,
				PasswordLastChangedAt:      timePointerFromNullTime(result.UserPasswordLastChangedAt),
				LastUpdatedAt:              timePointerFromNullTime(result.UserLastUpdatedAt),
				LastAcceptedTermsOfService: timePointerFromNullTime(result.UserLastAcceptedTermsOfService),
				LastAcceptedPrivacyPolicy:  timePointerFromNullTime(result.UserLastAcceptedPrivacyPolicy),
				TwoFactorSecretVerifiedAt:  timePointerFromNullTime(result.UserTwoFactorSecretVerifiedAt),
				AvatarSrc:                  stringPointerFromNullString(result.UserAvatarSrc),
				Birthday:                   timePointerFromNullTime(result.UserBirthday),
				ArchivedAt:                 timePointerFromNullTime(result.UserArchivedAt),
				AccountStatusExplanation:   result.UserUserAccountStatusExplanation,
				ID:                         result.UserID,
				AccountStatus:              result.UserUserAccountStatus,
				Username:                   result.UserUsername,
				FirstName:                  result.UserFirstName,
				LastName:                   result.UserLastName,
				EmailAddress:               result.UserEmailAddress,
				EmailAddressVerifiedAt:     timePointerFromNullTime(result.UserEmailAddressVerifiedAt),
				ServiceRole:                result.UserServiceRole,
				RequiresPasswordChange:     result.UserRequiresPasswordChange,
			},
			BelongsToHousehold: result.MembershipBelongsToHousehold,
			HouseholdRole:      result.MembershipHouseholdRole,
			DefaultHousehold:   result.MembershipDefaultHousehold,
		})
	}

	if household == nil {
		return nil, sql.ErrNoRows
	}

	return household, nil
}

// getHouseholdsForUser fetches a list of households from the database that meet a particular filter.
func (q *Querier) getHouseholdsForUser(ctx context.Context, querier database.SQLQueryExecutor, userID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.Household], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.Household]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetHouseholdsForUser(ctx, querier, &generated.GetHouseholdsForUserParams{
		BelongsToUser: userID,
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing households list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.Household{
			CreatedAt:                  result.CreatedAt,
			SubscriptionPlanID:         stringPointerFromNullString(result.SubscriptionPlanID),
			LastUpdatedAt:              timePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                 timePointerFromNullTime(result.ArchivedAt),
			ContactPhone:               result.ContactPhone,
			BillingStatus:              result.BillingStatus,
			AddressLine1:               result.AddressLine1,
			AddressLine2:               result.AddressLine2,
			City:                       result.City,
			State:                      result.State,
			ZipCode:                    result.ZipCode,
			Country:                    result.Country,
			Latitude:                   float64PointerFromNullString(result.Latitude),
			Longitude:                  float64PointerFromNullString(result.Longitude),
			PaymentProcessorCustomerID: result.PaymentProcessorCustomerID,
			BelongsToUser:              result.BelongsToUser,
			ID:                         result.ID,
			Name:                       result.Name,
			Members:                    nil,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetHouseholds fetches a list of households from the database that meet a particular filter.
func (q *Querier) GetHouseholds(ctx context.Context, userID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.Household], err error) {
	return q.getHouseholdsForUser(ctx, q.db, userID, filter)
}

// CreateHousehold creates a household in the database.
func (q *Querier) CreateHousehold(ctx context.Context, input *types.HouseholdDatabaseCreationInput) (*types.Household, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.UserIDKey, input.BelongsToUser)

	// begin household creation transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	// create the household.
	if writeErr := q.generatedQuerier.CreateHousehold(ctx, tx, &generated.CreateHouseholdParams{
		City:              input.City,
		Name:              input.Name,
		BillingStatus:     types.UnpaidHouseholdBillingStatus,
		ContactPhone:      input.ContactPhone,
		AddressLine1:      input.AddressLine1,
		AddressLine2:      input.AddressLine2,
		ID:                input.ID,
		State:             input.State,
		ZipCode:           input.ZipCode,
		Country:           input.Country,
		BelongsToUser:     input.BelongsToUser,
		WebhookHmacSecret: input.WebhookEncryptionKey,
		Latitude:          nullStringFromFloat64Pointer(input.Latitude),
		Longitude:         nullStringFromFloat64Pointer(input.Longitude),
	}); writeErr != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(writeErr, span, "creating household")
	}

	household := &types.Household{
		ID:            input.ID,
		Name:          input.Name,
		BelongsToUser: input.BelongsToUser,
		BillingStatus: types.UnpaidHouseholdBillingStatus,
		ContactPhone:  input.ContactPhone,
		AddressLine1:  input.AddressLine1,
		AddressLine2:  input.AddressLine2,
		City:          input.City,
		State:         input.State,
		ZipCode:       input.ZipCode,
		Country:       input.Country,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
		CreatedAt:     q.currentTime(),
	}

	if err = q.generatedQuerier.AddUserToHousehold(ctx, tx, &generated.AddUserToHouseholdParams{
		ID:                 identifiers.New(),
		BelongsToUser:      input.BelongsToUser,
		BelongsToHousehold: household.ID,
		HouseholdRole:      authorization.HouseholdAdminRole.String(),
	}); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing household membership creation query")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	tracing.AttachToSpan(span, keys.HouseholdIDKey, household.ID)
	logger.Info("household created")

	return household, nil
}

// UpdateHousehold updates a particular household. Note that UpdateHousehold expects the provided input to have a valid ID.
func (q *Querier) UpdateHousehold(ctx context.Context, updated *types.Household) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.HouseholdIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateHousehold(ctx, q.db, &generated.UpdateHouseholdParams{
		Name:          updated.Name,
		ContactPhone:  updated.ContactPhone,
		AddressLine1:  updated.AddressLine1,
		AddressLine2:  updated.AddressLine2,
		City:          updated.City,
		State:         updated.State,
		ZipCode:       updated.ZipCode,
		Country:       updated.Country,
		BelongsToUser: updated.BelongsToUser,
		ID:            updated.ID,
		Latitude:      nullStringFromFloat64Pointer(updated.Latitude),
		Longitude:     nullStringFromFloat64Pointer(updated.Longitude),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating household")
	}

	logger.Info("household updated")

	return nil
}

// ArchiveHousehold archives a household from the database by its ID.
func (q *Querier) ArchiveHousehold(ctx context.Context, householdID, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" || userID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.HouseholdIDKey, householdID)

	logger := q.logger.WithValues(map[string]any{
		keys.HouseholdIDKey: householdID,
		keys.UserIDKey:      userID,
	})

	if _, err := q.generatedQuerier.ArchiveHousehold(ctx, q.db, &generated.ArchiveHouseholdParams{
		BelongsToUser: userID,
		ID:            householdID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving household")
	}

	logger.Info("household archived")

	return nil
}
