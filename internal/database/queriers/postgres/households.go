package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	// householdsTableName is what the households table calls itself.
	householdsTableName = "households"

	householdUserMembershipsOnHouseholdsJoinClause = "household_user_memberships ON household_user_memberships.belongs_to_household = households.id"
	usersOnHouseholdUserMembershipsJoinClause      = "users ON household_user_memberships.belongs_to_user = users.id"
)

var (
	_ types.HouseholdDataManager = (*SQLQuerier)(nil)

	householdsTableColumns = []string{
		"households.id",
		"households.name",
		"households.billing_status",
		"households.contact_email",
		"households.contact_phone",
		"households.payment_processor_customer_id",
		"households.subscription_plan_id",
		"households.created_on",
		"households.last_updated_on",
		"households.archived_on",
		"households.belongs_to_user",
		"users.id",
		"users.username",
		"users.email_address",
		"users.avatar_src",
		"users.requires_password_change",
		"users.password_last_changed_on",
		"users.two_factor_secret_verified_on",
		"users.service_roles",
		"users.reputation",
		"users.reputation_explanation",
		"users.birth_day",
		"users.birth_month",
		"users.created_on",
		"users.last_updated_on",
		"users.archived_on",
	}
)

// scanHousehold takes a database Scanner (i.e. *sql.Row) and scans the result into a Household struct.
func (q *SQLQuerier) scanHousehold(ctx context.Context, scan database.Scanner, includeCounts bool) (household *types.Household, membership *types.HouseholdUserMembershipWithUser, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	household = &types.Household{Members: []*types.HouseholdUserMembershipWithUser{}}
	membership = &types.HouseholdUserMembershipWithUser{BelongsToUser: &types.User{}}

	var (
		rawHouseholdRoles,
		rawServiceRoles string
	)

	targetVars := []interface{}{
		&household.ID,
		&household.Name,
		&household.BillingStatus,
		&household.ContactEmail,
		&household.ContactPhone,
		&household.PaymentProcessorCustomerID,
		&household.SubscriptionPlanID,
		&household.CreatedOn,
		&household.LastUpdatedOn,
		&household.ArchivedOn,
		&household.BelongsToUser,
		&membership.BelongsToUser.ID,
		&membership.BelongsToUser.Username,
		&membership.BelongsToUser.EmailAddress,
		&membership.BelongsToUser.AvatarSrc,
		&membership.BelongsToUser.RequiresPasswordChange,
		&membership.BelongsToUser.PasswordLastChangedOn,
		&membership.BelongsToUser.TwoFactorSecretVerifiedOn,
		&rawServiceRoles,
		&membership.BelongsToUser.ServiceHouseholdStatus,
		&membership.BelongsToUser.ReputationExplanation,
		&membership.BelongsToUser.BirthDay,
		&membership.BelongsToUser.BirthMonth,
		&membership.BelongsToUser.CreatedOn,
		&membership.BelongsToUser.LastUpdatedOn,
		&membership.BelongsToUser.ArchivedOn,
		&membership.ID,
		&membership.BelongsToUser.ID,
		&membership.BelongsToHousehold,
		&rawHouseholdRoles,
		&membership.DefaultHousehold,
		&membership.CreatedOn,
		&membership.LastUpdatedOn,
		&membership.ArchivedOn,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, nil, 0, 0, observability.PrepareError(err, logger, span, "fetching memberships from database")
	}

	membership.HouseholdRoles = strings.Split(rawHouseholdRoles, householdMemberRolesSeparator)
	membership.BelongsToUser.ServiceRoles = strings.Split(rawServiceRoles, serviceRolesSeparator)

	return household, membership, filteredCount, totalCount, nil
}

// scanHouseholds takes some database rows and turns them into a slice of households.
func (q *SQLQuerier) scanHouseholds(ctx context.Context, rows database.ResultIterator, includeCounts bool) (households []*types.Household, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	households = []*types.Household{}

	var currentHousehold *types.Household
	for rows.Next() {
		household, membership, fc, tc, scanErr := q.scanHousehold(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if currentHousehold == nil {
			currentHousehold = household
		}

		if currentHousehold.ID != household.ID {
			households = append(households, currentHousehold)
			currentHousehold = household
		}

		currentHousehold.Members = append(currentHousehold.Members, membership)

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}
	}

	if currentHousehold != nil {
		households = append(households, currentHousehold)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return households, filteredCount, totalCount, nil
}

const getHouseholdQuery = `
	SELECT
		households.id,
		households.name,
		households.billing_status,
		households.contact_email,
		households.contact_phone,
		households.payment_processor_customer_id,
		households.subscription_plan_id,
		households.created_on,
		households.last_updated_on,
		households.archived_on,
		households.belongs_to_user,
        users.id,
        users.username,
        users.email_address,
        users.avatar_src,
        users.requires_password_change,
        users.password_last_changed_on,
        users.two_factor_secret_verified_on,
        users.service_roles,
        users.reputation,
        users.reputation_explanation,
        users.birth_day,
        users.birth_month,
        users.created_on,
        users.last_updated_on,
        users.archived_on,
		household_user_memberships.id,
		household_user_memberships.belongs_to_user,
		household_user_memberships.belongs_to_household,
		household_user_memberships.household_roles,
		household_user_memberships.default_household,
		household_user_memberships.created_on,
		household_user_memberships.last_updated_on,
		household_user_memberships.archived_on
	FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
	JOIN users ON household_user_memberships.belongs_to_user = users.id
	WHERE households.archived_on IS NULL
	AND households.belongs_to_user = $1
	AND households.id = $2
`

// GetHousehold fetches a household from the database.
func (q *SQLQuerier) GetHousehold(ctx context.Context, householdID, userID string) (*types.Household, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" || userID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachHouseholdIDToSpan(span, householdID)
	tracing.AttachUserIDToSpan(span, userID)

	logger := q.logger.WithValues(map[string]interface{}{
		keys.HouseholdIDKey: householdID,
		keys.UserIDKey:      userID,
	})

	args := []interface{}{
		userID,
		householdID,
	}

	rows, err := q.performReadQuery(ctx, q.db, "household", getHouseholdQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing households list retrieval query")
	}

	households, _, _, err := q.scanHouseholds(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	var household *types.Household
	if len(households) > 0 {
		household = households[0]
	}

	if household == nil {
		return nil, sql.ErrNoRows
	}

	return household, nil
}

const getHouseholdByIDQuery = `
	SELECT
		households.id,
		households.name,
		households.billing_status,
		households.contact_email,
		households.contact_phone,
		households.payment_processor_customer_id,
		households.subscription_plan_id,
		households.created_on,
		households.last_updated_on,
		households.archived_on,
		households.belongs_to_user,
        users.id,
        users.username,
        users.email_address,
        users.avatar_src,
        users.requires_password_change,
        users.password_last_changed_on,
        users.two_factor_secret_verified_on,
        users.service_roles,
        users.reputation,
        users.reputation_explanation,
        users.birth_day,
        users.birth_month,
        users.created_on,
        users.last_updated_on,
        users.archived_on,
		household_user_memberships.id,
		household_user_memberships.belongs_to_user,
		household_user_memberships.belongs_to_household,
		household_user_memberships.household_roles,
		household_user_memberships.default_household,
		household_user_memberships.created_on,
		household_user_memberships.last_updated_on,
		household_user_memberships.archived_on
	FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
	JOIN users ON household_user_memberships.belongs_to_user = users.id
	WHERE households.archived_on IS NULL
	AND households.id = $1
`

// GetHouseholdByID fetches a household from the database by its ID.
func (q *SQLQuerier) GetHouseholdByID(ctx context.Context, householdID string) (*types.Household, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachHouseholdIDToSpan(span, householdID)

	logger := q.logger.WithValue(keys.HouseholdIDKey, householdID)

	args := []interface{}{
		householdID,
	}

	rows, err := q.performReadQuery(ctx, q.db, "household", getHouseholdByIDQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing households list retrieval query")
	}

	households, _, _, err := q.scanHouseholds(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	var household *types.Household
	if len(households) > 0 {
		household = households[0]
	}

	if household == nil {
		return nil, sql.ErrNoRows
	}

	return household, nil
}

const getAllHouseholdsCountQuery = `
	SELECT COUNT(households.id) FROM households WHERE households.archived_on IS NULL
`

// GetAllHouseholdsCount fetches the count of households from the database that meet a particular filter.
func (q *SQLQuerier) GetAllHouseholdsCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getAllHouseholdsCountQuery, "fetching count of all households")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of households")
	}

	return count, nil
}

// buildGetHouseholdsQuery builds a SQL query selecting households that adhere to a given QueryFilter and belong to a given household,
// and returns both the query and the relevant args to pass to the query executor.
func (q *SQLQuerier) buildGetHouseholdsQuery(ctx context.Context, userID string, forAdmin bool, filter *types.QueryFilter) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)

	if filter != nil {
		tracing.AttachFilterToSpan(span, filter.Page, filter.Limit, string(filter.SortBy))
	}

	var includeArchived bool
	if filter != nil {
		includeArchived = filter.IncludeArchived
	}

	filteredCountQuery, filteredCountQueryArgs := q.buildFilteredCountQuery(ctx, householdsTableName, nil, nil, userOwnershipColumn, userID, forAdmin, includeArchived, filter)
	totalCountQuery, totalCountQueryArgs := q.buildTotalCountQuery(ctx, householdsTableName, nil, nil, userOwnershipColumn, userID, forAdmin, includeArchived)

	builder := q.sqlBuilder.Select(append(
		append(householdsTableColumns, householdsUserMembershipTableColumns...),
		fmt.Sprintf("(%s) as filtered_count", filteredCountQuery),
		fmt.Sprintf("(%s) as total_count", totalCountQuery),
	)...).
		From(householdsTableName).
		Join(householdUserMembershipsOnHouseholdsJoinClause).
		Join(usersOnHouseholdUserMembershipsJoinClause)

	if !forAdmin {
		where := squirrel.Eq{
			"households.archived_on": nil,
		}

		if userID != "" {
			where["household_user_memberships.belongs_to_user"] = userID
		}

		builder = builder.Where(where)
	}

	builder = builder.GroupBy(fmt.Sprintf(
		"%s.id, users.id, %s.id",
		householdsTableName,
		householdsUserMembershipTableName,
	))

	if filter != nil {
		builder = applyFilterToQueryBuilder(filter, householdsTableName, builder)
	}

	query, selectArgs := q.buildQuery(span, builder)

	return query, append(append(filteredCountQueryArgs, totalCountQueryArgs...), selectArgs...)
}

// getHouseholds fetches a list of households from the database that meet a particular filter.
func (q *SQLQuerier) getHouseholds(ctx context.Context, querier database.SQLQueryExecutor, userID string, forAdmin bool, filter *types.QueryFilter) (x *types.HouseholdList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" && !forAdmin {
		return nil, ErrInvalidIDProvided
	}

	logger := filter.AttachToLogger(q.logger).WithValue(keys.UserIDKey, userID).WithValue("filter_is_nil", filter == nil)
	tracing.AttachQueryFilterToSpan(span, filter)
	tracing.AttachUserIDToSpan(span, userID)

	x = &types.HouseholdList{}
	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildGetHouseholdsQuery(ctx, userID, forAdmin, filter)

	rows, err := q.performReadQuery(ctx, querier, "households", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing households list retrieval query")
	}

	if x.Households, x.FilteredCount, x.TotalCount, err = q.scanHouseholds(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning households from database")
	}

	return x, nil
}

// GetHouseholds fetches a list of households from the database that meet a particular filter.
func (q *SQLQuerier) GetHouseholds(ctx context.Context, userID string, filter *types.QueryFilter) (x *types.HouseholdList, err error) {
	return q.getHouseholds(ctx, q.db, userID, false, filter)
}

// GetHouseholdsForAdmin fetches a list of households from the database that meet a particular filter for all users.
func (q *SQLQuerier) GetHouseholdsForAdmin(ctx context.Context, userID string, filter *types.QueryFilter) (x *types.HouseholdList, err error) {
	return q.getHouseholds(ctx, q.db, userID, true, filter)
}

const householdCreationQuery = `
	INSERT INTO households (id,name,billing_status,contact_email,contact_phone,belongs_to_user) VALUES ($1,$2,$3,$4,$5,$6)
`

const addUserToHouseholdDuringCreationQuery = `
	INSERT INTO household_user_memberships (id,belongs_to_user,belongs_to_household,household_roles)
	VALUES ($1,$2,$3,$4)
`

// CreateHousehold creates a household in the database.
func (q *SQLQuerier) CreateHousehold(ctx context.Context, input *types.HouseholdDatabaseCreationInput) (*types.Household, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.UserIDKey, input.BelongsToUser)

	// begin household creation transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "beginning transaction")
	}

	householdCreationArgs := []interface{}{
		input.ID,
		input.Name,
		types.UnpaidHouseholdBillingStatus,
		input.ContactEmail,
		input.ContactPhone,
		input.BelongsToUser,
	}

	// create the household.
	if writeErr := q.performWriteQuery(ctx, tx, "household creation", householdCreationQuery, householdCreationArgs); writeErr != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(writeErr, logger, span, "creating household")
	}

	household := &types.Household{
		ID:            input.ID,
		Name:          input.Name,
		BelongsToUser: input.BelongsToUser,
		BillingStatus: types.UnpaidHouseholdBillingStatus,
		ContactEmail:  input.ContactEmail,
		ContactPhone:  input.ContactPhone,
		CreatedOn:     q.currentTime(),
	}

	addInput := &types.HouseholdUserMembershipCreationRequestInput{
		ID:             ksuid.New().String(),
		UserID:         input.BelongsToUser,
		HouseholdID:    household.ID,
		HouseholdRoles: []string{authorization.HouseholdAdminRole.String()},
	}

	addUserToHouseholdArgs := []interface{}{
		addInput.ID,
		addInput.UserID,
		addInput.HouseholdID,
		strings.Join(addInput.HouseholdRoles, householdMemberRolesSeparator),
	}

	if err = q.performWriteQuery(ctx, tx, "household user membership creation", addUserToHouseholdDuringCreationQuery, addUserToHouseholdArgs); err != nil {
		q.rollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, logger, span, "performing household membership creation query")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareError(err, logger, span, "committing transaction")
	}

	tracing.AttachHouseholdIDToSpan(span, household.ID)
	logger.Info("household created")

	return household, nil
}

const updateHouseholdQuery = `
	UPDATE households SET name = $1, contact_email = $2, contact_phone = $3, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $4 AND id = $5
`

// UpdateHousehold updates a particular household. Note that UpdateHousehold expects the provided input to have a valid ID.
func (q *SQLQuerier) UpdateHousehold(ctx context.Context, updated *types.Household) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.HouseholdIDKey, updated.ID)
	tracing.AttachHouseholdIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Name,
		updated.ContactEmail,
		updated.ContactPhone,
		updated.BelongsToUser,
		updated.ID,
	}

	logger.WithValue("query", updateHouseholdQuery).WithValue("args", args).Info("making query for households")

	if err := q.performWriteQuery(ctx, q.db, "household update", updateHouseholdQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating household")
	}

	logger.Info("household updated")

	return nil
}

const archiveHouseholdQuery = `
	UPDATE households SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2
`

// ArchiveHousehold archives a household from the database by its ID.
func (q *SQLQuerier) ArchiveHousehold(ctx context.Context, householdID, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" || userID == "" {
		return ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	logger := q.logger.WithValues(map[string]interface{}{
		keys.HouseholdIDKey: householdID,
		keys.UserIDKey:      userID,
	})

	args := []interface{}{
		userID,
		householdID,
	}

	if err := q.performWriteQuery(ctx, q.db, "household archive", archiveHouseholdQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "archiving household")
	}

	logger.Info("household archived")

	return nil
}
