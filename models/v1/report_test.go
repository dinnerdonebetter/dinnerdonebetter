package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestReport_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &Report{}

		expected := &ReportUpdateInput{
			ReportType: fake.Word(),
			Concern:    fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.ReportType, i.ReportType)
		assert.Equal(t, expected.Concern, i.Concern)
	})
}

func TestReport_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		report := &Report{
			ReportType: fake.Word(),
			Concern:    fake.Word(),
		}

		expected := &ReportUpdateInput{
			ReportType: report.ReportType,
			Concern:    report.Concern,
		}
		actual := report.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
