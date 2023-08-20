package postgres

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
)

func timePointerFromNullTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}

	return nil
}

func stringPointerFromNullString(nt sql.NullString) *string {
	if nt.Valid {
		return &nt.String
	}

	return nil
}

func nullStringFromString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func nullStringFromStringPointer(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}

	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func nullTimeFromTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func nullTimeFromTimePointer(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}

func nullInt32FromUint8Pointer(i *uint8) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Int32: int32(*i),
		Valid: true,
	}
}

func nullInt32FromUint16(i uint16) sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(i),
		Valid: true,
	}
}

func boolFromNullBool(b sql.NullBool) bool {
	if b.Valid {
		return b.Bool
	}

	return false
}

func nullInt32FromInt32Pointer(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Int32: *i,
		Valid: true,
	}
}

func int32PointerFromNullInt32(i sql.NullInt32) *int32 {
	if i.Valid {
		return &i.Int32
	}

	return nil
}

func nullBoolFromBool(b bool) sql.NullBool {
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}

func nullFloat64FromFloat32Pointer(f *float32) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}

	return sql.NullFloat64{
		Float64: float64(*f),
		Valid:   true,
	}
}

func float32PointerFromNullString(f sql.NullString) *float32 {
	if f.Valid {
		if parsedFloat, err := strconv.ParseFloat(f.String, 64); err == nil {
			return pointers.Pointer(float32(parsedFloat))
		}
	}

	return nil
}
