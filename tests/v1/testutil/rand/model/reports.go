package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomReportCreationInput creates a random ReportInput
func RandomReportCreationInput() *models.ReportCreationInput {
	x := &models.ReportCreationInput{
		ReportType: fake.Word(),
		Concern:    fake.Word(),
	}

	return x
}
