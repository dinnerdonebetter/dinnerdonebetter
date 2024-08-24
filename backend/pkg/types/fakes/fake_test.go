package fakes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_buildFakeNumber(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotZero(t, buildFakeNumber())
	})
}
