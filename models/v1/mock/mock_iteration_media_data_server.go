package mock

import (
	"net/http"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.IterationMediaDataServer = (*IterationMediaDataServer)(nil)

// IterationMediaDataServer is a mocked models.IterationMediaDataServer for testing.
type IterationMediaDataServer struct {
	mock.Mock
}

// CreationInputMiddleware implements our interface requirements.
func (m *IterationMediaDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// UpdateInputMiddleware implements our interface requirements.
func (m *IterationMediaDataServer) UpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// ListHandler implements our interface requirements.
func (m *IterationMediaDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// CreateHandler implements our interface requirements.
func (m *IterationMediaDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ExistenceHandler implements our interface requirements.
func (m *IterationMediaDataServer) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ReadHandler implements our interface requirements.
func (m *IterationMediaDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// UpdateHandler implements our interface requirements.
func (m *IterationMediaDataServer) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ArchiveHandler implements our interface requirements.
func (m *IterationMediaDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
