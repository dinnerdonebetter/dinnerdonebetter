package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomReport retrieves a random report from the list of available reports
func fetchRandomReport(c *client.V1Client) *models.Report {
	reportsRes, err := c.GetReports(context.Background(), nil)
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
				return c.BuildCreateReportRequest(context.Background(), randmodel.RandomReportCreationInput())
			},
			Weight: 100,
		},
		"GetReport": {
			Name: "GetReport",
			Action: func() (*http.Request, error) {
				if randomReport := fetchRandomReport(c); randomReport != nil {
					return c.BuildGetReportRequest(context.Background(), randomReport.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetReports": {
			Name: "GetReports",
			Action: func() (*http.Request, error) {
				return c.BuildGetReportsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateReport": {
			Name: "UpdateReport",
			Action: func() (*http.Request, error) {
				if randomReport := fetchRandomReport(c); randomReport != nil {
					randomReport.ReportType = randmodel.RandomReportCreationInput().ReportType
					randomReport.Concern = randmodel.RandomReportCreationInput().Concern
					return c.BuildUpdateReportRequest(context.Background(), randomReport)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveReport": {
			Name: "ArchiveReport",
			Action: func() (*http.Request, error) {
				if randomReport := fetchRandomReport(c); randomReport != nil {
					return c.BuildArchiveReportRequest(context.Background(), randomReport.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
