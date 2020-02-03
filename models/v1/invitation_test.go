package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvitation_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &Invitation{}

		expected := &InvitationUpdateInput{
			Code:     "example",
			Consumed: false,
		}

		i.Update(expected)
		assert.Equal(t, expected.Code, i.Code)
		assert.Equal(t, expected.Consumed, i.Consumed)
	})
}
