package mock

import (
	"net/http"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.ValidIngredientDataServer = (*ValidIngredientDataServer)(nil)

// ValidIngredientDataServer is a mocked models.ValidIngredientDataServer for testing.
type ValidIngredientDataServer struct {
	mock.Mock
}

// CreationInputMiddleware implements our interface requirements.
func (m *ValidIngredientDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// UpdateInputMiddleware implements our interface requirements.
func (m *ValidIngredientDataServer) UpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// SearchHandler implements our interface requirements.
func (m *ValidIngredientDataServer) SearchHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ListHandler implements our interface requirements.
func (m *ValidIngredientDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// CreateHandler implements our interface requirements.
func (m *ValidIngredientDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ExistenceHandler implements our interface requirements.
func (m *ValidIngredientDataServer) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ReadHandler implements our interface requirements.
func (m *ValidIngredientDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// UpdateHandler implements our interface requirements.
func (m *ValidIngredientDataServer) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ArchiveHandler implements our interface requirements.
func (m *ValidIngredientDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
