package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertFloat32RangeWithOptionalMaxUpdateRequestInputToFloat32RangeWithOptionalMax(input *messages.Float32RangeWithOptionalMaxUpdateRequestInput) *messages.Float32RangeWithOptionalMax {

output := &messages.Float32RangeWithOptionalMax{
    Max: input.Max,
    Min: input.Min,
}

return output
}


func ConvertUint16RangeWithOptionalMaxUpdateRequestInputToUint16RangeWithOptionalMax(input *messages.Uint16RangeWithOptionalMaxUpdateRequestInput) *messages.Uint16RangeWithOptionalMax {

output := &messages.Uint16RangeWithOptionalMax{
    Min: input.Min,
    Max: input.Max,
}

return output
}

func ConvertUint32RangeWithOptionalMaxUpdateRequestInputToUint32RangeWithOptionalMax(input *messages.Uint32RangeWithOptionalMaxUpdateRequestInput) *messages.Uint32RangeWithOptionalMax {

output := &messages.Uint32RangeWithOptionalMax{
    Max: input.Max,
    Min: input.Min,
}

return output
}


