package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createOAuth2ClientTokenForTest(t *testing.T, ctx context.Context, exampleOAuth2ClientToken *types.OAuth2ClientToken, dbc *Querier) *types.OAuth2ClientToken {
	t.Helper()

	// create
	if exampleOAuth2ClientToken == nil {
		user := createUserForTest(t, ctx, nil, dbc)
		oauth2Client := createOAuth2ClientForTest(t, ctx, nil, dbc)
		exampleOAuth2ClientToken = fakes.BuildFakeOAuth2ClientToken()
		exampleOAuth2ClientToken.BelongsToUser = user.ID
		exampleOAuth2ClientToken.ClientID = oauth2Client.ID
	}
	dbInput := converters.ConvertOAuth2ClientTokenToOAuth2ClientTokenDatabaseCreationInput(exampleOAuth2ClientToken)

	created, err := dbc.CreateOAuth2ClientToken(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleOAuth2ClientToken.CodeCreatedAt = created.CodeCreatedAt
	exampleOAuth2ClientToken.AccessCreatedAt = created.AccessCreatedAt
	exampleOAuth2ClientToken.RefreshCreatedAt = created.RefreshCreatedAt
	assert.Equal(t, exampleOAuth2ClientToken, created)

	oauth2ClientToken, err := dbc.GetOAuth2ClientTokenByAccess(ctx, created.Access)
	assert.NoError(t, err)
	exampleOAuth2ClientToken.CodeCreatedAt = oauth2ClientToken.CodeCreatedAt
	exampleOAuth2ClientToken.AccessCreatedAt = oauth2ClientToken.AccessCreatedAt
	exampleOAuth2ClientToken.RefreshCreatedAt = oauth2ClientToken.RefreshCreatedAt
	assert.Equal(t, oauth2ClientToken, exampleOAuth2ClientToken)

	return created
}

func TestQuerier_Integration_OAuth2ClientTokens(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)
	oauth2Client := createOAuth2ClientForTest(t, ctx, nil, dbc)
	exampleOAuth2ClientToken := fakes.BuildFakeOAuth2ClientToken()
	exampleOAuth2ClientToken.BelongsToUser = user.ID
	exampleOAuth2ClientToken.ClientID = oauth2Client.ClientID
	exampleOAuth2ClientToken.Scope = "household_member"

	// create
	createdOAuth2ClientToken := createOAuth2ClientTokenForTest(t, ctx, exampleOAuth2ClientToken, dbc)

	// delete
	assert.NoError(t, dbc.ArchiveOAuth2ClientTokenByCode(ctx, createdOAuth2ClientToken.Code))

	// create
	createdOAuth2ClientToken = createOAuth2ClientTokenForTest(t, ctx, exampleOAuth2ClientToken, dbc)

	// delete
	assert.NoError(t, dbc.ArchiveOAuth2ClientTokenByAccess(ctx, createdOAuth2ClientToken.Access))

	// create
	createdOAuth2ClientToken = createOAuth2ClientTokenForTest(t, ctx, exampleOAuth2ClientToken, dbc)

	// delete
	assert.NoError(t, dbc.ArchiveOAuth2ClientTokenByRefresh(ctx, createdOAuth2ClientToken.Refresh))
}