package mock

import (
	"net/http"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.PreparationDataServer = (*PreparationDataServer)(nil)

// PreparationDataServer is a mocked models.PreparationDataServer for testing
type PreparationDataServer struct {
	mock.Mock
}

// CreationInputMiddleware implements our interface requirements
func (m *PreparationDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// UpdateInputMiddleware implements our interface requirements
func (m *PreparationDataServer) UpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// ListHandler implements our interface requirements
func (m *PreparationDataServer) ListHandler() http.HandlerFunc {
	args := m.Called()
	return args.Get(0).(http.HandlerFunc)
}

// CreateHandler implements our interface requirements
func (m *PreparationDataServer) CreateHandler() http.HandlerFunc {
	args := m.Called()
	return args.Get(0).(http.HandlerFunc)
}

// ReadHandler implements our interface requirements
func (m *PreparationDataServer) ReadHandler() http.HandlerFunc {
	args := m.Called()
	return args.Get(0).(http.HandlerFunc)
}

// UpdateHandler implements our interface requirements
func (m *PreparationDataServer) UpdateHandler() http.HandlerFunc {
	args := m.Called()
	return args.Get(0).(http.HandlerFunc)
}

// ArchiveHandler implements our interface requirements
func (m *PreparationDataServer) ArchiveHandler() http.HandlerFunc {
	args := m.Called()
	return args.Get(0).(http.HandlerFunc)
}
