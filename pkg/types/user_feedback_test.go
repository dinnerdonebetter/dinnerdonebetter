package types

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserFeedbackCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &UserFeedbackCreationRequestInput{
			Prompt:   t.Name(),
			Feedback: t.Name(),
			Rating:   1.23,
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestUserFeedbackDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &UserFeedbackDatabaseCreationInput{
			ID:       t.Name(),
			Prompt:   t.Name(),
			Feedback: t.Name(),
			Rating:   1.23,
			ByUser:   t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
