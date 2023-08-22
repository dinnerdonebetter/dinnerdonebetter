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

func stringFromNullString(nt sql.NullString) string {
	if nt.Valid {
		return nt.String
	}

	return ""
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

func nullInt32FromUint16Pointer(i *uint16) sql.NullInt32 {
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

func nullInt32FromUint32Pointer(i *uint32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Int32: int32(*i),
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

func float64PointerFromNullString(f sql.NullString) *float64 {
	if f.Valid {
		if parsedFloat, err := strconv.ParseFloat(f.String, 64); err == nil {
			return &parsedFloat
		}
	}

	return nil
}

func stringFromFloat32(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

func float32FromString(s string) float32 {
	if parsedFloat, err := strconv.ParseFloat(s, 64); err == nil {
		return float32(parsedFloat)
	}

	return 0
}

func float32FromNullString(s sql.NullString) float32 {
	if s.Valid {
		return float32FromString(s.String)
	}

	return 0
}

func nullStringFromFloat32Pointer(f *float32) sql.NullString {
	if f == nil {
		return sql.NullString{}
	}

	return sql.NullString{
		String: stringFromFloat32(*f),
		Valid:  true,
	}
}

func nullStringFromFloat32(f float32) sql.NullString {
	return sql.NullString{
		String: stringFromFloat32(f),
		Valid:  true,
	}
}

func stringFromFloat64(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func nullStringFromFloat64Pointer(f *float64) sql.NullString {
	if f == nil {
		return sql.NullString{}
	}

	return sql.NullString{
		String: stringFromFloat64(*f),
		Valid:  true,
	}
}

func nullInt64FromUint32Pointer(f *uint32) sql.NullInt64 {
	if f == nil {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: int64(*f),
		Valid: true,
	}
}

func uint16PointerFromNullInt32(f sql.NullInt32) *uint16 {
	if f.Valid {
		return pointers.Pointer(uint16(f.Int32))
	}

	return nil
}

func uint32PointerFromNullInt32(f sql.NullInt32) *uint32 {
	if f.Valid {
		return pointers.Pointer(uint32(f.Int32))
	}

	return nil
}

func uint32PointerFromNullInt64(f sql.NullInt64) *uint32 {
	if f.Valid {
		return pointers.Pointer(uint32(f.Int64))
	}

	return nil
}
