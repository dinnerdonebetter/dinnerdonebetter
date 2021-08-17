package frontend

import (
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProvideService(t *testing.T) {
	t.Parallel()

	cfg := &Config{}
	logger := logging.NewNoopLogger()
	authService := &mocktypes.AuthService{}

	s := ProvideService(cfg, logger, authService)

	mock.AssertExpectationsForObjects(t, authService)
	assert.NotNil(t, s)
}
