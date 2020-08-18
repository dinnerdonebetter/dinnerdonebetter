package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeReport builds a faked report.
func BuildFakeReport() *models.Report {
	return &models.Report{
		ID:            fake.Uint64(),
		ReportType:    fake.Word(),
		Concern:       fake.Word(),
		CreatedOn:     uint64(uint32(fake.Date().Unix())),
		BelongsToUser: fake.Uint64(),
	}
}

// BuildFakeReportList builds a faked ReportList.
func BuildFakeReportList() *models.ReportList {
	exampleReport1 := BuildFakeReport()
	exampleReport2 := BuildFakeReport()
	exampleReport3 := BuildFakeReport()

	return &models.ReportList{
		Pagination: models.Pagination{
			Page:  1,
			Limit: 20,
		},
		Reports: []models.Report{
			*exampleReport1,
			*exampleReport2,
			*exampleReport3,
		},
	}
}

// BuildFakeReportUpdateInputFromReport builds a faked ReportUpdateInput from a report.
func BuildFakeReportUpdateInputFromReport(report *models.Report) *models.ReportUpdateInput {
	return &models.ReportUpdateInput{
		ReportType:    report.ReportType,
		Concern:       report.Concern,
		BelongsToUser: report.BelongsToUser,
	}
}

// BuildFakeReportCreationInput builds a faked ReportCreationInput.
func BuildFakeReportCreationInput() *models.ReportCreationInput {
	report := BuildFakeReport()
	return BuildFakeReportCreationInputFromReport(report)
}

// BuildFakeReportCreationInputFromReport builds a faked ReportCreationInput from a report.
func BuildFakeReportCreationInputFromReport(report *models.Report) *models.ReportCreationInput {
	return &models.ReportCreationInput{
		ReportType:    report.ReportType,
		Concern:       report.Concern,
		BelongsToUser: report.BelongsToUser,
	}
}
