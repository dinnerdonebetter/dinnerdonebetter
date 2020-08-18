package mock

import (
	"net/http"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.InvitationDataServer = (*InvitationDataServer)(nil)

// InvitationDataServer is a mocked models.InvitationDataServer for testing.
type InvitationDataServer struct {
	mock.Mock
}

// CreationInputMiddleware implements our interface requirements.
func (m *InvitationDataServer) CreationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// UpdateInputMiddleware implements our interface requirements.
func (m *InvitationDataServer) UpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// ListHandler implements our interface requirements.
func (m *InvitationDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// CreateHandler implements our interface requirements.
func (m *InvitationDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ExistenceHandler implements our interface requirements.
func (m *InvitationDataServer) ExistenceHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ReadHandler implements our interface requirements.
func (m *InvitationDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// UpdateHandler implements our interface requirements.
func (m *InvitationDataServer) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ArchiveHandler implements our interface requirements.
func (m *InvitationDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
