package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type methodHolder struct{}

func (m *methodHolder) Method() {}

func TestGetMethodName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &methodHolder{}
		assert.Equal(t, "Method", GetMethodName(x.Method))
	})
}
