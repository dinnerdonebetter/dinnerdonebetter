package frontend

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_fetchInvitation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitation",
			testutils.ContextMatcher,
			exampleInvitation.ID,
		).Return(exampleInvitation, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		actual, err := s.service.fetchInvitation(s.ctx, req)
		assert.Equal(t, exampleInvitation, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching invitation", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitation",
			testutils.ContextMatcher,
			exampleInvitation.ID,
		).Return((*types.Invitation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		actual, err := s.service.fetchInvitation(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachInvitationCreationInputToRequest(input *types.InvitationCreationInput) *http.Request {
	form := url.Values{
		invitationCreationInputCodeFormKey:     {anyToString(input.Code)},
		invitationCreationInputConsumedFormKey: {anyToString(input.Consumed)},
	}

	return httptest.NewRequest(http.MethodPost, "/invitations", strings.NewReader(form.Encode()))
}

func TestService_buildInvitationCreatorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationCreatorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationCreatorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationCreatorView(false)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationCreatorView(true)(res, req)
	})

	T.Run("without base template and error writing to response", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		res := &testutils.MockHTTPResponseWriter{}
		res.On("Write", mock.Anything).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationCreatorView(false)(res, req)
	})
}

func TestService_parseFormEncodedInvitationCreationInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		expected := fakes.BuildFakeInvitationCreationInput()
		expected.BelongsToAccount = s.exampleAccount.ID
		req := attachInvitationCreationInputToRequest(expected)

		actual := s.service.parseFormEncodedInvitationCreationInput(s.ctx, req, s.sessionCtxData)
		assert.Equal(t, expected, actual)
	})

	T.Run("with error extracting form from request", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedInvitationCreationInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.InvitationCreationInput{}
		req := attachInvitationCreationInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedInvitationCreationInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})
}

func TestService_handleInvitationCreationRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		exampleInput := fakes.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
		exampleInput.BelongsToAccount = s.sessionCtxData.ActiveAccountID

		res := httptest.NewRecorder()
		req := attachInvitationCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"CreateInvitation",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleInvitation, nil)
		s.service.dataStore = mockDB

		s.service.handleInvitationCreationRequest(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.NotEmpty(t, res.Header().Get(htmxRedirectionHeader))

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		exampleInput := fakes.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachInvitationCreationInputToRequest(exampleInput)

		s.service.handleInvitationCreationRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		exampleInput := fakes.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
		exampleInput.BelongsToAccount = s.sessionCtxData.ActiveAccountID

		res := httptest.NewRecorder()
		req := attachInvitationCreationInputToRequest(&types.InvitationCreationInput{})

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"CreateInvitation",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return(exampleInvitation, nil)
		s.service.dataStore = mockDB

		s.service.handleInvitationCreationRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error creating invitation in database", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		exampleInput := fakes.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)
		exampleInput.BelongsToAccount = s.sessionCtxData.ActiveAccountID

		res := httptest.NewRecorder()
		req := attachInvitationCreationInputToRequest(exampleInput)

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"CreateInvitation",
			testutils.ContextMatcher,
			exampleInput,
			s.sessionCtxData.Requester.UserID,
		).Return((*types.Invitation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		s.service.handleInvitationCreationRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildInvitationEditorView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitation",
			testutils.ContextMatcher,
			exampleInvitation.ID,
		).Return(exampleInvitation, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationEditorView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitation",
			testutils.ContextMatcher,
			exampleInvitation.ID,
		).Return(exampleInvitation, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationEditorView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationEditorView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching invitation", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitation",
			testutils.ContextMatcher,
			exampleInvitation.ID,
		).Return((*types.Invitation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationEditorView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_fetchInvitations(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitationList := fakes.BuildFakeInvitationList()

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleInvitationList, nil)
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		actual, err := s.service.fetchInvitations(s.ctx, req)
		assert.Equal(t, exampleInvitationList, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.InvitationList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		actual, err := s.service.fetchInvitations(s.ctx, req)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_buildInvitationsTableView(T *testing.T) {
	T.Parallel()

	T.Run("with base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitationList := fakes.BuildFakeInvitationList()
		for _, invitation := range exampleInvitationList.Invitations {
			invitation.BelongsToAccount = s.exampleAccount.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleInvitationList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationsTableView(true)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("without base template", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitationList := fakes.BuildFakeInvitationList()

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleInvitationList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationsTableView(false)(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)
		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationsTableView(true)(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.InvitationList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/invitations", nil)

		s.service.buildInvitationsTableView(true)(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func attachInvitationUpdateInputToRequest(input *types.InvitationUpdateInput) *http.Request {
	form := url.Values{
		invitationUpdateInputCodeFormKey:     {anyToString(input.Code)},
		invitationUpdateInputConsumedFormKey: {anyToString(input.Consumed)},
	}

	return httptest.NewRequest(http.MethodPost, "/invitations", strings.NewReader(form.Encode()))
}

func TestService_parseFormEncodedInvitationUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		expected := fakes.BuildFakeInvitationUpdateInputFromInvitation(exampleInvitation)

		req := attachInvitationUpdateInputToRequest(expected)

		actual := s.service.parseFormEncodedInvitationUpdateInput(s.ctx, req, s.sessionCtxData)
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		badBody := &testutils.MockReadCloser{}
		badBody.On("Read", mock.IsType([]byte{})).Return(0, errors.New("blah"))

		req := httptest.NewRequest(http.MethodGet, "/test", badBody)

		actual := s.service.parseFormEncodedInvitationUpdateInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input attached to valid form", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.InvitationUpdateInput{}

		req := attachInvitationUpdateInputToRequest(exampleInput)

		actual := s.service.parseFormEncodedInvitationUpdateInput(s.ctx, req, s.sessionCtxData)
		assert.Nil(t, actual)
	})
}

func TestService_handleInvitationUpdateRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		exampleInput := fakes.BuildFakeInvitationUpdateInputFromInvitation(exampleInvitation)

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitation",
			testutils.ContextMatcher,
			exampleInvitation.ID,
		).Return(exampleInvitation, nil)

		mockDB.InvitationDataManager.On(
			"UpdateInvitation",
			testutils.ContextMatcher,
			exampleInvitation,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachInvitationUpdateInputToRequest(exampleInput)

		s.service.handleInvitationUpdateRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		exampleInput := fakes.BuildFakeInvitationUpdateInputFromInvitation(exampleInvitation)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := attachInvitationUpdateInputToRequest(exampleInput)

		s.service.handleInvitationUpdateRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInput := &types.InvitationUpdateInput{}

		res := httptest.NewRecorder()
		req := attachInvitationUpdateInputToRequest(exampleInput)

		s.service.handleInvitationUpdateRequest(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		exampleInput := fakes.BuildFakeInvitationUpdateInputFromInvitation(exampleInvitation)

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitation",
			testutils.ContextMatcher,
			exampleInvitation.ID,
		).Return((*types.Invitation)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachInvitationUpdateInputToRequest(exampleInput)

		s.service.handleInvitationUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error updating data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		exampleInput := fakes.BuildFakeInvitationUpdateInputFromInvitation(exampleInvitation)

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"GetInvitation",
			testutils.ContextMatcher,
			exampleInvitation.ID,
		).Return(exampleInvitation, nil)

		mockDB.InvitationDataManager.On(
			"UpdateInvitation",
			testutils.ContextMatcher,
			exampleInvitation,
			s.sessionCtxData.Requester.UserID,
			[]*types.FieldChangeSummary(nil),
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := attachInvitationUpdateInputToRequest(exampleInput)

		s.service.handleInvitationUpdateRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestService_handleInvitationArchiveRequest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		exampleInvitationList := fakes.BuildFakeInvitationList()

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"ArchiveInvitation",
			testutils.ContextMatcher,
			exampleInvitation.ID,
			s.sessionCtxData.ActiveAccountID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.InvitationDataManager.On(
			"GetInvitations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleInvitationList, nil)
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/invitations", nil)

		s.service.handleInvitationArchiveRequest(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		s.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/invitations", nil)

		s.service.handleInvitationArchiveRequest(res, req)

		assert.Equal(t, unauthorizedRedirectResponseCode, res.Code)
	})

	T.Run("with error archiving invitation", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"ArchiveInvitation",
			testutils.ContextMatcher,
			exampleInvitation.ID,
			s.sessionCtxData.ActiveAccountID,
			s.sessionCtxData.Requester.UserID,
		).Return(errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/invitations", nil)

		s.service.handleInvitationArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error retrieving new list of invitations", func(t *testing.T) {
		t.Parallel()

		s := buildTestHelper(t)

		exampleInvitation := fakes.BuildFakeInvitation()
		exampleInvitation.BelongsToAccount = s.exampleAccount.ID
		s.service.invitationIDFetcher = func(*http.Request) uint64 {
			return exampleInvitation.ID
		}

		mockDB := database.BuildMockDatabase()
		mockDB.InvitationDataManager.On(
			"ArchiveInvitation",
			testutils.ContextMatcher,
			exampleInvitation.ID,
			s.sessionCtxData.ActiveAccountID,
			s.sessionCtxData.Requester.UserID,
		).Return(nil)
		s.service.dataStore = mockDB

		mockDB.InvitationDataManager.On(
			"GetInvitations",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.InvitationList)(nil), errors.New("blah"))
		s.service.dataStore = mockDB

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/invitations", nil)

		s.service.handleInvitationArchiveRequest(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
