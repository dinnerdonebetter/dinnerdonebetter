package mock

import (
	"net/http"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.ValidInstrumentDataServer = (*ValidInstrumentDataServer)(nil)

// ValidInstrumentDataServer is a mocked models.ValidInstrumentDataServer for testing.
type ValidInstrumentDataServer struct {
	mock.Mock
}

// CreationInputMiddleware implements our interface requirements.
func (m *ValidInstrumentDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// UpdateInputMiddleware implements our interface requirements.
func (m *ValidInstrumentDataServer) UpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// SearchHandler implements our interface requirements.
func (m *ValidInstrumentDataServer) SearchHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ListHandler implements our interface requirements.
func (m *ValidInstrumentDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// CreateHandler implements our interface requirements.
func (m *ValidInstrumentDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ExistenceHandler implements our interface requirements.
func (m *ValidInstrumentDataServer) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ReadHandler implements our interface requirements.
func (m *ValidInstrumentDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// UpdateHandler implements our interface requirements.
func (m *ValidInstrumentDataServer) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ArchiveHandler implements our interface requirements.
func (m *ValidInstrumentDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
