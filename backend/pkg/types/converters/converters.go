package converters

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ConvertTimestampToPBTimestamp(t time.Time) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}

func ConvertTimestampPointerToPBTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return &timestamppb.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}

func float32PointerToFloat64Pointer(f *float32) *float64 {
	if f == nil {
		return nil
	}

	x := float64(*f)
	return &x
}
