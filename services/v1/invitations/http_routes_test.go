package invitations

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockencoding "gitlab.com/prixfixe/prixfixe/internal/v1/encoding/mock"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/v1/metrics/mock"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"
	mockmodels "gitlab.com/prixfixe/prixfixe/models/v1/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocknewsman "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
)

func TestInvitationsService_ListHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitationList := fakemodels.BuildFakeInvitationList()

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleInvitationList, nil)
		s.invitationDataManager = invitationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.InvitationList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.InvitationList)(nil), sql.ErrNoRows)
		s.invitationDataManager = invitationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.InvitationList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager, ed)
	})

	T.Run("with error fetching invitations from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.InvitationList)(nil), errors.New("blah"))
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitationList := fakemodels.BuildFakeInvitationList()

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleInvitationList, nil)
		s.invitationDataManager = invitationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.InvitationList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager, ed)
	})
}

func TestInvitationsService_CreateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("CreateInvitation", mock.Anything, mock.AnythingOfType("*models.InvitationCreationInput")).Return(exampleInvitation, nil)
		s.invitationDataManager = invitationDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.invitationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Invitation")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager, mc, r, ed)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error creating invitation", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("CreateInvitation", mock.Anything, mock.AnythingOfType("*models.InvitationCreationInput")).Return(exampleInvitation, errors.New("blah"))
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("CreateInvitation", mock.Anything, mock.AnythingOfType("*models.InvitationCreationInput")).Return(exampleInvitation, nil)
		s.invitationDataManager = invitationDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.invitationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Invitation")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), CreateMiddlewareCtxKey, exampleInput))

		s.CreateHandler()(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager, mc, r, ed)
	})
}

func TestInvitationsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("InvitationExists", mock.Anything, exampleInvitation.ID).Return(true, nil)
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})

	T.Run("with no such invitation in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("InvitationExists", mock.Anything, exampleInvitation.ID).Return(false, sql.ErrNoRows)
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})

	T.Run("with error fetching invitation from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("InvitationExists", mock.Anything, exampleInvitation.ID).Return(false, errors.New("blah"))
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})
}

func TestInvitationsService_ReadHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitation", mock.Anything, exampleInvitation.ID).Return(exampleInvitation, nil)
		s.invitationDataManager = invitationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Invitation")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager, ed)
	})

	T.Run("with no such invitation in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitation", mock.Anything, exampleInvitation.ID).Return((*models.Invitation)(nil), sql.ErrNoRows)
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})

	T.Run("with error fetching invitation from database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitation", mock.Anything, exampleInvitation.ID).Return((*models.Invitation)(nil), errors.New("blah"))
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitation", mock.Anything, exampleInvitation.ID).Return(exampleInvitation, nil)
		s.invitationDataManager = invitationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Invitation")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager, ed)
	})
}

func TestInvitationsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationUpdateInputFromInvitation(exampleInvitation)

		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitation", mock.Anything, exampleInvitation.ID).Return(exampleInvitation, nil)
		invitationDataManager.On("UpdateInvitation", mock.Anything, mock.AnythingOfType("*models.Invitation")).Return(nil)
		s.invitationDataManager = invitationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Invitation")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, invitationDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with no rows fetching invitation", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationUpdateInputFromInvitation(exampleInvitation)

		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitation", mock.Anything, exampleInvitation.ID).Return((*models.Invitation)(nil), sql.ErrNoRows)
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})

	T.Run("with error fetching invitation", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationUpdateInputFromInvitation(exampleInvitation)

		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitation", mock.Anything, exampleInvitation.ID).Return((*models.Invitation)(nil), errors.New("blah"))
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})

	T.Run("with error updating invitation", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationUpdateInputFromInvitation(exampleInvitation)

		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitation", mock.Anything, exampleInvitation.ID).Return(exampleInvitation, nil)
		invitationDataManager.On("UpdateInvitation", mock.Anything, mock.AnythingOfType("*models.Invitation")).Return(errors.New("blah"))
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationUpdateInputFromInvitation(exampleInvitation)

		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("GetInvitation", mock.Anything, exampleInvitation.ID).Return(exampleInvitation, nil)
		invitationDataManager.On("UpdateInvitation", mock.Anything, mock.AnythingOfType("*models.Invitation")).Return(nil)
		s.invitationDataManager = invitationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.Invitation")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), UpdateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler()(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, invitationDataManager, ed)
	})
}

func TestInvitationsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	userIDFetcher := func(_ *http.Request) uint64 {
		return exampleUser.ID
	}

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("ArchiveInvitation", mock.Anything, exampleInvitation.ID, exampleUser.ID).Return(nil)
		s.invitationDataManager = invitationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.invitationCounter = mc

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager, mc, r)
	})

	T.Run("with no invitation in database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("ArchiveInvitation", mock.Anything, exampleInvitation.ID, exampleUser.ID).Return(sql.ErrNoRows)
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		s.userIDFetcher = userIDFetcher

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		s.invitationIDFetcher = func(req *http.Request) uint64 {
			return exampleInvitation.ID
		}

		invitationDataManager := &mockmodels.InvitationDataManager{}
		invitationDataManager.On("ArchiveInvitation", mock.Anything, exampleInvitation.ID, exampleUser.ID).Return(errors.New("blah"))
		s.invitationDataManager = invitationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://prixfixe.app",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, invitationDataManager)
	})
}
