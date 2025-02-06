package fake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type example struct {
	Name string
	Age  int
}

func TestBuildFakeForTest(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := BuildFakeForTest[*example](t)

		assert.NotEmpty(t, actual.Name)
		assert.NotEmpty(t, actual.Age)
		assert.NotNil(t, actual)
	})
}
