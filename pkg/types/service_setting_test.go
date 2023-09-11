package types

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v6"
)

func TestServiceSetting_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ServiceSetting{}
		input := &ServiceSettingUpdateRequestInput{}

		fake.Struct(&input)
		input.AdminsOnly = pointers.Pointer(true)

		x.Update(input)
	})
}
