package converters

import (
	"time"

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
