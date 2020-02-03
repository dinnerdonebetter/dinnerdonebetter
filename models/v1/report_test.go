package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &Report{}

		expected := &ReportUpdateInput{
			ReportType: "example",
			Concern:    "example",
		}

		i.Update(expected)
		assert.Equal(t, expected.ReportType, i.ReportType)
		assert.Equal(t, expected.Concern, i.Concern)
	})
}
