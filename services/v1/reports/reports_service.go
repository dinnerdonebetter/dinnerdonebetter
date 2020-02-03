package reports

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// CreateMiddlewareCtxKey is a string alias we can use for referring to report input data in contexts
	CreateMiddlewareCtxKey models.ContextKey = "report_create_input"
	// UpdateMiddlewareCtxKey is a string alias we can use for referring to report update data in contexts
	UpdateMiddlewareCtxKey models.ContextKey = "report_update_input"

	counterName        metrics.CounterName = "reports"
	counterDescription                     = "the number of reports managed by the reports service"
	topicName          string              = "reports"
	serviceName        string              = "reports_service"
)

var (
	_ models.ReportDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list reports
	Service struct {
		logger          logging.Logger
		reportCounter   metrics.UnitCounter
		reportDatabase  models.ReportDataManager
		userIDFetcher   UserIDFetcher
		reportIDFetcher ReportIDFetcher
		encoderDecoder  encoding.EncoderDecoder
		reporter        newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs
	UserIDFetcher func(*http.Request) uint64

	// ReportIDFetcher is a function that fetches report IDs
	ReportIDFetcher func(*http.Request) uint64
)

// ProvideReportsService builds a new ReportsService
func ProvideReportsService(
	ctx context.Context,
	logger logging.Logger,
	db models.ReportDataManager,
	userIDFetcher UserIDFetcher,
	reportIDFetcher ReportIDFetcher,
	encoder encoding.EncoderDecoder,
	reportCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	reportCounter, err := reportCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:          logger.WithName(serviceName),
		reportDatabase:  db,
		encoderDecoder:  encoder,
		reportCounter:   reportCounter,
		userIDFetcher:   userIDFetcher,
		reportIDFetcher: reportIDFetcher,
		reporter:        reporter,
	}

	reportCount, err := svc.reportDatabase.GetAllReportsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("setting current report count: %w", err)
	}
	svc.reportCounter.IncrementBy(ctx, reportCount)

	return svc, nil
}
