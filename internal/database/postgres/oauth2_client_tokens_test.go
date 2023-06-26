package postgres

import (
	"context"
	"testing"

	encryptionmock "github.com/dinnerdonebetter/backend/internal/pkg/cryptography/mock"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQuerier_ArchiveOAuth2ClientTokenByAccess(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2ClientToken := fakes.BuildFakeOAuth2ClientToken()
		exampleEncryptedAccess := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		encDec := encryptionmock.NewMockEncryptorDecryptor()
		encDec.On("Encrypt", testutils.ContextMatcher, exampleOAuth2ClientToken.Access).Return(exampleEncryptedAccess, nil)
		c.oauth2ClientTokenEncDec = encDec

		args := []any{
			exampleEncryptedAccess,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveOAuth2ClientTokenByAccessQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.ArchiveOAuth2ClientTokenByAccess(ctx, exampleOAuth2ClientToken.Access)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveOAuth2ClientTokenByCode(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2ClientToken := fakes.BuildFakeOAuth2ClientToken()
		exampleEncryptedCode := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		encDec := encryptionmock.NewMockEncryptorDecryptor()
		encDec.On("Encrypt", testutils.ContextMatcher, exampleOAuth2ClientToken.Code).Return(exampleEncryptedCode, nil)
		c.oauth2ClientTokenEncDec = encDec

		args := []any{
			exampleEncryptedCode,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveOAuth2ClientTokenByCodeQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.ArchiveOAuth2ClientTokenByCode(ctx, exampleOAuth2ClientToken.Code)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveOAuth2ClientTokenByRefresh(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2ClientToken := fakes.BuildFakeOAuth2ClientToken()
		exampleEncryptedRefresh := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		encDec := encryptionmock.NewMockEncryptorDecryptor()
		encDec.On("Encrypt", testutils.ContextMatcher, exampleOAuth2ClientToken.Refresh).Return(exampleEncryptedRefresh, nil)
		c.oauth2ClientTokenEncDec = encDec

		args := []any{
			exampleEncryptedRefresh,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveOAuth2ClientTokenByRefreshQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.ArchiveOAuth2ClientTokenByRefresh(ctx, exampleOAuth2ClientToken.Refresh)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}
