package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"
)

// fetchRandomReport retrieves a random report from the list of available reports.
func fetchRandomReport(ctx context.Context, c *client.V1Client) *models.Report {
	reportsRes, err := c.GetReports(ctx, nil)
	if err != nil || reportsRes == nil || len(reportsRes.Reports) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(reportsRes.Reports))
	return &reportsRes.Reports[randIndex]
}

func buildReportActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateReport": {
			Name: "CreateReport",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				reportInput := fakemodels.BuildFakeReportCreationInput()

				return c.BuildCreateReportRequest(ctx, reportInput)
			},
			Weight: 100,
		},
		"GetReport": {
			Name: "GetReport",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomReport := fetchRandomReport(ctx, c)
				if randomReport == nil {
					return nil, fmt.Errorf("retrieving random report: %w", ErrUnavailableYet)
				}

				return c.BuildGetReportRequest(ctx, randomReport.ID)
			},
			Weight: 100,
		},
		"GetReports": {
			Name: "GetReports",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				return c.BuildGetReportsRequest(ctx, nil)
			},
			Weight: 100,
		},
		"UpdateReport": {
			Name: "UpdateReport",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				if randomReport := fetchRandomReport(ctx, c); randomReport != nil {
					newReport := fakemodels.BuildFakeReportCreationInput()
					randomReport.ReportType = newReport.ReportType
					randomReport.Concern = newReport.Concern
					return c.BuildUpdateReportRequest(ctx, randomReport)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveReport": {
			Name: "ArchiveReport",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomReport := fetchRandomReport(ctx, c)
				if randomReport == nil {
					return nil, fmt.Errorf("retrieving random report: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveReportRequest(ctx, randomReport.ID)
			},
			Weight: 85,
		},
	}
}
