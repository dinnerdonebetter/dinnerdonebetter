package converters

import (
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/pointer"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertTimeToPBTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func ConvertTimePointerToPBTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func ConvertPBTimestampToTime(t *timestamppb.Timestamp) time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.AsTime()
}

func ConvertPBTimestampToTimePointer(t *timestamppb.Timestamp) *time.Time {
	return pointer.To(ConvertPBTimestampToTime(t))
}
