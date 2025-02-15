package fakes

import (
	"fmt"
	"math"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types"

	fake "github.com/brianvoe/gofakeit/v7"
)

func init() {
	if err := fake.Seed(time.Now().UnixNano()); err != nil {
		panic(err)
	}
}

const (
	exampleQuantity = 3
)

// BuildFakeID builds a fake ID.
func BuildFakeID() string {
	return identifiers.New()
}

func BuildFakeString() string {
	return fake.Word()
}

func buildFakeTOTPToken() string {
	return fmt.Sprintf("%d%s", fake.Number(0, 9), fake.Zip())
}

// BuildFakeTime builds a fake time.
func BuildFakeTime() time.Time {
	return fake.Date().Add(0).Truncate(time.Second).UTC()
}

func buildFakeNumber() float64 {
	return math.Round(float64((fake.Number(101, math.MaxInt8-1) * 100) / 100))
}

// buildUniqueString builds a fake string.
func buildUniqueString() string {
	return fake.LoremIpsumSentence(7)
}

// buildFakePassword builds a fake password.
func buildFakePassword() string {
	return fake.Password(true, true, true, true, false, 32)
}

func BuildFakeFloat32RangeWithOptionalMax() types.Float32RangeWithOptionalMax {
	minimum := float32(buildFakeNumber())
	return types.Float32RangeWithOptionalMax{
		Min: minimum,
		Max: pointer.To(float32(buildFakeNumber()) + minimum),
	}
}

func BuildFakeOptionalFloat32Range() types.OptionalFloat32Range {
	minimum := float32(buildFakeNumber())
	return types.OptionalFloat32Range{
		Min: pointer.To(minimum),
		Max: pointer.To(float32(buildFakeNumber()) + minimum),
	}
}

func BuildFakeOptionalUint32Range() types.OptionalUint32Range {
	minimum := uint32(buildFakeNumber())
	return types.OptionalUint32Range{
		Min: pointer.To(minimum),
		Max: pointer.To(uint32(buildFakeNumber()) + minimum),
	}
}

func BuildFakeUint16RangeWithOptionalMax() types.Uint16RangeWithOptionalMax {
	minimum := uint16(buildFakeNumber())
	return types.Uint16RangeWithOptionalMax{
		Min: minimum,
		Max: pointer.To(uint16(buildFakeNumber()) + minimum),
	}
}

func BuildFakeUint32RangeWithOptionalMax() types.Uint32RangeWithOptionalMax {
	minimum := uint32(buildFakeNumber())
	return types.Uint32RangeWithOptionalMax{
		Min: minimum,
		Max: pointer.To(uint32(buildFakeNumber()) + minimum),
	}
}

func BuildFakeUint32RangeWithOptionalMaxUpdateRequestInput() types.Uint32RangeWithOptionalMaxUpdateRequestInput {
	minimum := uint32(buildFakeNumber())
	return types.Uint32RangeWithOptionalMaxUpdateRequestInput{
		Min: &minimum,
		Max: pointer.To(uint32(buildFakeNumber()) + minimum),
	}
}
