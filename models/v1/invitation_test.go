package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestInvitation_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &Invitation{}

		expected := &InvitationUpdateInput{
			Code:     fake.Word(),
			Consumed: fake.Bool(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Code, i.Code)
		assert.Equal(t, expected.Consumed, i.Consumed)
	})
}

func TestInvitation_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		invitation := &Invitation{
			Code:     fake.Word(),
			Consumed: fake.Bool(),
		}

		expected := &InvitationUpdateInput{
			Code:     invitation.Code,
			Consumed: invitation.Consumed,
		}
		actual := invitation.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
