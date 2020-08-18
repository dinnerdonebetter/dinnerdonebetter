package models

import (
	"context"
	"net/http"
)

type (
	// Report represents a report.
	Report struct {
		ID            uint64  `json:"id"`
		ReportType    string  `json:"reportType"`
		Concern       string  `json:"concern"`
		CreatedOn     uint64  `json:"createdOn"`
		LastUpdatedOn *uint64 `json:"lastUpdatedOn"`
		ArchivedOn    *uint64 `json:"archivedOn"`
		BelongsToUser uint64  `json:"belongsToUser"`
	}

	// ReportList represents a list of reports.
	ReportList struct {
		Pagination
		Reports []Report `json:"reports"`
	}

	// ReportCreationInput represents what a user could set as input for creating reports.
	ReportCreationInput struct {
		ReportType    string `json:"reportType"`
		Concern       string `json:"concern"`
		BelongsToUser uint64 `json:"-"`
	}

	// ReportUpdateInput represents what a user could set as input for updating reports.
	ReportUpdateInput struct {
		ReportType    string `json:"reportType"`
		Concern       string `json:"concern"`
		BelongsToUser uint64 `json:"-"`
	}

	// ReportDataManager describes a structure capable of storing reports permanently.
	ReportDataManager interface {
		ReportExists(ctx context.Context, reportID uint64) (bool, error)
		GetReport(ctx context.Context, reportID uint64) (*Report, error)
		GetAllReportsCount(ctx context.Context) (uint64, error)
		GetAllReports(ctx context.Context, resultChannel chan []Report) error
		GetReports(ctx context.Context, filter *QueryFilter) (*ReportList, error)
		GetReportsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]Report, error)
		CreateReport(ctx context.Context, input *ReportCreationInput) (*Report, error)
		UpdateReport(ctx context.Context, updated *Report) error
		ArchiveReport(ctx context.Context, reportID, userID uint64) error
	}

	// ReportDataServer describes a structure capable of serving traffic related to reports.
	ReportDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ReportInput with a report.
func (x *Report) Update(input *ReportUpdateInput) {
	if input.ReportType != "" && input.ReportType != x.ReportType {
		x.ReportType = input.ReportType
	}

	if input.Concern != "" && input.Concern != x.Concern {
		x.Concern = input.Concern
	}
}

// ToUpdateInput creates a ReportUpdateInput struct for a report.
func (x *Report) ToUpdateInput() *ReportUpdateInput {
	return &ReportUpdateInput{
		ReportType: x.ReportType,
		Concern:    x.Concern,
	}
}
