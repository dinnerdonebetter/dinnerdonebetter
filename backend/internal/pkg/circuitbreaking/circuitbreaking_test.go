package circuitbreaking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvideCircuitBreaker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, ProvideCircuitBreaker(nil))
	})
}

func TestEnsureCircuitBreaker(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, EnsureCircuitBreaker(nil))
}
