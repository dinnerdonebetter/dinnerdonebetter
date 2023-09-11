package types

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v6"
)

func TestServiceSettingConfiguration_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ServiceSettingConfiguration{}
		input := &ServiceSettingConfigurationUpdateRequestInput{}

		fake.Struct(&input)

		x.Update(input)
	})
}
