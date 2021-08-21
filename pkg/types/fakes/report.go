package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeReport builds a faked report.
func BuildFakeReport() *types.Report {
	return &types.Report{
		ID:                 uint64(fake.Uint32()),
		ExternalID:         fake.UUID(),
		ReportType:         fake.Word(),
		Concern:            fake.Word(),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
		BelongsToHousehold: fake.Uint64(),
	}
}

// BuildFakeReportList builds a faked ReportList.
func BuildFakeReportList() *types.ReportList {
	var examples []*types.Report
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeReport())
	}

	return &types.ReportList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Reports: examples,
	}
}

// BuildFakeReportUpdateInput builds a faked ReportUpdateInput from a report.
func BuildFakeReportUpdateInput() *types.ReportUpdateInput {
	report := BuildFakeReport()
	return &types.ReportUpdateInput{
		ReportType:         report.ReportType,
		Concern:            report.Concern,
		BelongsToHousehold: report.BelongsToHousehold,
	}
}

// BuildFakeReportUpdateInputFromReport builds a faked ReportUpdateInput from a report.
func BuildFakeReportUpdateInputFromReport(report *types.Report) *types.ReportUpdateInput {
	return &types.ReportUpdateInput{
		ReportType:         report.ReportType,
		Concern:            report.Concern,
		BelongsToHousehold: report.BelongsToHousehold,
	}
}

// BuildFakeReportCreationInput builds a faked ReportCreationInput.
func BuildFakeReportCreationInput() *types.ReportCreationInput {
	report := BuildFakeReport()
	return BuildFakeReportCreationInputFromReport(report)
}

// BuildFakeReportCreationInputFromReport builds a faked ReportCreationInput from a report.
func BuildFakeReportCreationInputFromReport(report *types.Report) *types.ReportCreationInput {
	return &types.ReportCreationInput{
		ReportType:         report.ReportType,
		Concern:            report.Concern,
		BelongsToHousehold: report.BelongsToHousehold,
	}
}
