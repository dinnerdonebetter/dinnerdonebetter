package mock

import (
	"net/http"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.ValidPreparationDataServer = (*ValidPreparationDataServer)(nil)

// ValidPreparationDataServer is a mocked models.ValidPreparationDataServer for testing.
type ValidPreparationDataServer struct {
	mock.Mock
}

// CreationInputMiddleware implements our interface requirements.
func (m *ValidPreparationDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// UpdateInputMiddleware implements our interface requirements.
func (m *ValidPreparationDataServer) UpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// SearchHandler implements our interface requirements.
func (m *ValidPreparationDataServer) SearchHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ListHandler implements our interface requirements.
func (m *ValidPreparationDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// CreateHandler implements our interface requirements.
func (m *ValidPreparationDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ExistenceHandler implements our interface requirements.
func (m *ValidPreparationDataServer) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ReadHandler implements our interface requirements.
func (m *ValidPreparationDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// UpdateHandler implements our interface requirements.
func (m *ValidPreparationDataServer) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ArchiveHandler implements our interface requirements.
func (m *ValidPreparationDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
