package mock

import (
	"net/http"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.RecipeStepInstrumentDataServer = (*RecipeStepInstrumentDataServer)(nil)

// RecipeStepInstrumentDataServer is a mocked models.RecipeStepInstrumentDataServer for testing.
type RecipeStepInstrumentDataServer struct {
	mock.Mock
}

// CreationInputMiddleware implements our interface requirements.
func (m *RecipeStepInstrumentDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// UpdateInputMiddleware implements our interface requirements.
func (m *RecipeStepInstrumentDataServer) UpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// ListHandler implements our interface requirements.
func (m *RecipeStepInstrumentDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// CreateHandler implements our interface requirements.
func (m *RecipeStepInstrumentDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ExistenceHandler implements our interface requirements.
func (m *RecipeStepInstrumentDataServer) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ReadHandler implements our interface requirements.
func (m *RecipeStepInstrumentDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// UpdateHandler implements our interface requirements.
func (m *RecipeStepInstrumentDataServer) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ArchiveHandler implements our interface requirements.
func (m *RecipeStepInstrumentDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
