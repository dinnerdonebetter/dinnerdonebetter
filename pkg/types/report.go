package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Report represents a report.
	Report struct {
		ArchivedOn       *uint64 `json:"archivedOn"`
		LastUpdatedOn    *uint64 `json:"lastUpdatedOn"`
		ExternalID       string  `json:"externalID"`
		ReportType       string  `json:"reportType"`
		Concern          string  `json:"concern"`
		CreatedOn        uint64  `json:"createdOn"`
		ID               uint64  `json:"id"`
		BelongsToAccount uint64  `json:"belongsToAccount"`
	}

	// ReportList represents a list of reports.
	ReportList struct {
		Reports []*Report `json:"reports"`
		Pagination
	}

	// ReportCreationInput represents what a user could set as input for creating reports.
	ReportCreationInput struct {
		ReportType       string `json:"reportType"`
		Concern          string `json:"concern"`
		BelongsToAccount uint64 `json:"-"`
	}

	// ReportUpdateInput represents what a user could set as input for updating reports.
	ReportUpdateInput struct {
		ReportType       string `json:"reportType"`
		Concern          string `json:"concern"`
		BelongsToAccount uint64 `json:"-"`
	}

	// ReportDataManager describes a structure capable of storing reports permanently.
	ReportDataManager interface {
		ReportExists(ctx context.Context, reportID uint64) (bool, error)
		GetReport(ctx context.Context, reportID uint64) (*Report, error)
		GetAllReportsCount(ctx context.Context) (uint64, error)
		GetAllReports(ctx context.Context, resultChannel chan []*Report, bucketSize uint16) error
		GetReports(ctx context.Context, filter *QueryFilter) (*ReportList, error)
		GetReportsWithIDs(ctx context.Context, accountID uint64, limit uint8, ids []uint64) ([]*Report, error)
		CreateReport(ctx context.Context, input *ReportCreationInput, createdByUser uint64) (*Report, error)
		UpdateReport(ctx context.Context, updated *Report, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveReport(ctx context.Context, reportID, accountID, archivedBy uint64) error
		GetAuditLogEntriesForReport(ctx context.Context, reportID uint64) ([]*AuditLogEntry, error)
	}

	// ReportDataService describes a structure capable of serving traffic related to reports.
	ReportDataService interface {
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ReportUpdateInput with a report.
func (x *Report) Update(input *ReportUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.ReportType != x.ReportType {
		out = append(out, &FieldChangeSummary{
			FieldName: "ReportType",
			OldValue:  x.ReportType,
			NewValue:  input.ReportType,
		})

		x.ReportType = input.ReportType
	}

	if input.Concern != x.Concern {
		out = append(out, &FieldChangeSummary{
			FieldName: "Concern",
			OldValue:  x.Concern,
			NewValue:  input.Concern,
		})

		x.Concern = input.Concern
	}

	return out
}

var _ validation.ValidatableWithContext = (*ReportCreationInput)(nil)

// ValidateWithContext validates a ReportCreationInput.
func (x *ReportCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ReportType, validation.Required),
		validation.Field(&x.Concern, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*ReportUpdateInput)(nil)

// ValidateWithContext validates a ReportUpdateInput.
func (x *ReportUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ReportType, validation.Required),
		validation.Field(&x.Concern, validation.Required),
	)
}
