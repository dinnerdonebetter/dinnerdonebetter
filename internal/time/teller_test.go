package time

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandardTimeTeller_Now(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, (&StandardTimeTeller{}).Now())
	})
}
