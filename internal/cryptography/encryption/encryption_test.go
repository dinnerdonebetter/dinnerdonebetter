package encryption

import (
	"fmt"
	"github.com/prixfixeco/backend/pkg/types/fakes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncrypt(T *testing.T) {
	T.Parallel()

	T.Run("basic operation", func(t *testing.T) {
		t.Parallel()

		expected := t.Name()
		user := fakes.BuildFakeUser()
		secret := fmt.Sprintf("%s.%s", user.ID, user.CreatedAt.Format("15:04:05.00"))

		encrypted, err := Encrypt(expected, secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, encrypted)

		actual, err := Decrypt(encrypted, secret)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
