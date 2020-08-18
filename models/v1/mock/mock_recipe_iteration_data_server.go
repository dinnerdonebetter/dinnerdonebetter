package mock

import (
	"net/http"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeIterationDataServer = (*RecipeIterationDataServer)(nil)

// RecipeIterationDataServer is a mocked models.RecipeIterationDataServer for testing.
type RecipeIterationDataServer struct {
	mock.Mock
}

// CreationInputMiddleware implements our interface requirements.
func (m *RecipeIterationDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// UpdateInputMiddleware implements our interface requirements.
func (m *RecipeIterationDataServer) UpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// ListHandler implements our interface requirements.
func (m *RecipeIterationDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// CreateHandler implements our interface requirements.
func (m *RecipeIterationDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ExistenceHandler implements our interface requirements.
func (m *RecipeIterationDataServer) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ReadHandler implements our interface requirements.
func (m *RecipeIterationDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// UpdateHandler implements our interface requirements.
func (m *RecipeIterationDataServer) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ArchiveHandler implements our interface requirements.
func (m *RecipeIterationDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
