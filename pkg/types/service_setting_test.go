package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestServiceSetting_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ServiceSetting{}
		input := &ServiceSettingUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.AdminsOnly = pointers.Pointer(true)

		x.Update(input)
	})
}

func TestServiceSettingCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ServiceSettingCreationRequestInput{
			Name:        t.Name(),
			Description: t.Name(),
			Type:        t.Name(),
			Enumeration: []string{
				t.Name(),
			},
			DefaultValue: pointers.Pointer(t.Name()),
		}

		actual := x.ValidateWithContext(ctx)

		assert.NoError(t, actual)
	})

	T.Run("with invalid struct", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ServiceSettingCreationRequestInput{}

		actual := x.ValidateWithContext(ctx)

		assert.Error(t, actual)
	})

	T.Run("with invalid default value", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ServiceSettingCreationRequestInput{
			Name:        t.Name(),
			Description: t.Name(),
			Type:        t.Name(),
			Enumeration: []string{
				t.Name(),
			},
			DefaultValue: pointers.Pointer("whatever"),
		}

		actual := x.ValidateWithContext(ctx)

		assert.Error(t, actual)
	})
}

func TestServiceSettingDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ServiceSettingDatabaseCreationInput{
			ID:          t.Name(),
			Name:        t.Name(),
			Description: t.Name(),
			Type:        t.Name(),
			Enumeration: []string{
				t.Name(),
			},
			DefaultValue: pointers.Pointer(t.Name()),
		}

		actual := x.ValidateWithContext(ctx)

		assert.NoError(t, actual)
	})
}

func TestServiceSettingUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ServiceSettingUpdateRequestInput{
			Name:        pointers.Pointer(t.Name()),
			Description: pointers.Pointer(t.Name()),
			Enumeration: []string{
				t.Name(),
			},
			DefaultValue: pointers.Pointer(t.Name()),
		}

		actual := x.ValidateWithContext(ctx)

		assert.NoError(t, actual)
	})

	T.Run("with invalid struct", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ServiceSettingUpdateRequestInput{}

		actual := x.ValidateWithContext(ctx)

		assert.Error(t, actual)
	})

	T.Run("with incorrect default value", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ServiceSettingUpdateRequestInput{
			Name:        pointers.Pointer(t.Name()),
			Description: pointers.Pointer(t.Name()),
			Enumeration: []string{
				t.Name(),
			},
			DefaultValue: pointers.Pointer("whatever"),
		}

		actual := x.ValidateWithContext(ctx)

		assert.Error(t, actual)
	})
}
