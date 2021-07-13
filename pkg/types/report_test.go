package types

import (
	"context"
	"encoding/json"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReport_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Report{}

		updated := &ReportUpdateInput{
			ReportType: fake.Word(),
			Concern:    fake.Word(),
		}

		expected := []*FieldChangeSummary{
			{
				FieldName: "ReportType",
				OldValue:  x.ReportType,
				NewValue:  updated.ReportType,
			},
			{
				FieldName: "Concern",
				OldValue:  x.Concern,
				NewValue:  updated.Concern,
			},
		}
		actual := x.Update(updated)

		expectedJSONBytes, err := json.Marshal(expected)
		require.NoError(t, err)

		actualJSONBytes, err := json.Marshal(actual)
		require.NoError(t, err)

		expectedJSON, actualJSON := string(expectedJSONBytes), string(actualJSONBytes)

		assert.Equal(t, expectedJSON, actualJSON)

		assert.Equal(t, updated.ReportType, x.ReportType)
		assert.Equal(t, updated.Concern, x.Concern)
	})
}

func TestReportCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ReportCreationInput{
			ReportType: fake.Word(),
			Concern:    fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ReportCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestReportUpdateInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ReportUpdateInput{
			ReportType: fake.Word(),
			Concern:    fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ReportUpdateInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
