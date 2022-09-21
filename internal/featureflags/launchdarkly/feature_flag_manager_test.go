package launchdarkly

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"

	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestFeatureFlagManager_CanUseFeature(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUsername := fakes.BuildFakeUser().Username

		ffm := &FeatureFlagManager{logger: logging.NewNoopLogger(), tracer: tracing.NewTracerForTest(t.Name())}

		fakeClient := &mockClient{}
		user := lduser.NewUserBuilder(exampleUsername).Name(exampleUsername).Build()
		fakeClient.On("BoolVariation", t.Name(), user, false).Return(true, nil)
		ffm.client = fakeClient

		actual, err := ffm.CanUseFeature(ctx, exampleUsername, t.Name())
		assert.NoError(t, err)
		assert.True(t, actual)
	})
}
