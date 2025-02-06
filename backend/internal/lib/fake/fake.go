package fake

import (
	"math"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/interfaces"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/stretchr/testify/require"
)

// BuildFakeID builds a fake ID.
func BuildFakeID() string {
	return identifiers.New()
}

// BuildFakeForTest builds a fake instance of insert-struct-here for a test.
func BuildFakeForTest[X any](t *testing.T) (x X) {
	t.Helper()
	require.NoError(t, faker.FakeData(&x,
		options.WithRandomFloatBoundaries(interfaces.RandomFloatBoundary{Start: 1.0, End: float64(math.MaxUint8)}),
		options.WithRandomIntegerBoundaries(interfaces.RandomIntegerBoundary{Start: 1, End: math.MaxUint8}),
	))
	return x
}
