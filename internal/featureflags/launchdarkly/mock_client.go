package launchdarkly

import (
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

func (m *mockClient) BoolVariation(key string, context ldcontext.Context, defaultVal bool) (bool, error) {
	args := m.Called(key, context, defaultVal)
	return args.Bool(0), args.Error(1)
}
