package identifiers

import (
	"testing"

	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := New()
		assert.NotEmpty(t, actual)
	})
}

func TestValidate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := Validate(ksuid.New().String())
		assert.NoError(t, actual)
	})
}

func Test_newID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := newID(false)
		assert.NotEmpty(t, actual)
	})

	T.Run("with xid", func(t *testing.T) {
		t.Parallel()

		actual := newID(true)
		assert.NotEmpty(t, actual)
	})
}

func Test_parseID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := parseID(ksuid.New().String(), false)
		assert.NoError(t, actual)
	})

	T.Run("with xid", func(t *testing.T) {
		t.Parallel()

		actual := parseID(xid.New().String(), true)
		assert.NoError(t, actual)
	})
}
