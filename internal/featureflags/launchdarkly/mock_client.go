package launchdarkly

import (
	"github.com/stretchr/testify/mock"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
)

type mockClient struct {
	mock.Mock
}

func (m *mockClient) BoolVariation(key string, user lduser.User, defaultVal bool) (bool, error) {
	args := m.Called(key, user, defaultVal)
	return args.Bool(0), args.Error(1)
}
