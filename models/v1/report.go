package models

import (
	"context"
	"net/http"
)

type (
	// Report represents a report
	Report struct {
		ID         uint64  `json:"id"`
		ReportType string  `json:"report_type"`
		Concern    string  `json:"concern"`
		CreatedOn  uint64  `json:"created_on"`
		UpdatedOn  *uint64 `json:"updated_on"`
		ArchivedOn *uint64 `json:"archived_on"`
		BelongsTo  uint64  `json:"belongs_to"`
	}

	// ReportList represents a list of reports
	ReportList struct {
		Pagination
		Reports []Report `json:"reports"`
	}

	// ReportCreationInput represents what a user could set as input for creating reports
	ReportCreationInput struct {
		ReportType string `json:"report_type"`
		Concern    string `json:"concern"`
		BelongsTo  uint64 `json:"-"`
	}

	// ReportUpdateInput represents what a user could set as input for updating reports
	ReportUpdateInput struct {
		ReportType string `json:"report_type"`
		Concern    string `json:"concern"`
		BelongsTo  uint64 `json:"-"`
	}

	// ReportDataManager describes a structure capable of storing reports permanently
	ReportDataManager interface {
		GetReport(ctx context.Context, reportID, userID uint64) (*Report, error)
		GetReportCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllReportsCount(ctx context.Context) (uint64, error)
		GetReports(ctx context.Context, filter *QueryFilter, userID uint64) (*ReportList, error)
		GetAllReportsForUser(ctx context.Context, userID uint64) ([]Report, error)
		CreateReport(ctx context.Context, input *ReportCreationInput) (*Report, error)
		UpdateReport(ctx context.Context, updated *Report) error
		ArchiveReport(ctx context.Context, id, userID uint64) error
	}

	// ReportDataServer describes a structure capable of serving traffic related to reports
	ReportDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an ReportInput with a report
func (x *Report) Update(input *ReportUpdateInput) {
	if input.ReportType != "" && input.ReportType != x.ReportType {
		x.ReportType = input.ReportType
	}

	if input.Concern != "" && input.Concern != x.Concern {
		x.Concern = input.Concern
	}
}

// ToInput creates a ReportUpdateInput struct for a report
func (x *Report) ToInput() *ReportUpdateInput {
	return &ReportUpdateInput{
		ReportType: x.ReportType,
		Concern:    x.Concern,
	}
}
