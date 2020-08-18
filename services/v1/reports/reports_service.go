package reports

import (
	"fmt"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

const (
	// createMiddlewareCtxKey is a string alias we can use for referring to report input data in contexts.
	createMiddlewareCtxKey models.ContextKey = "report_create_input"
	// updateMiddlewareCtxKey is a string alias we can use for referring to report update data in contexts.
	updateMiddlewareCtxKey models.ContextKey = "report_update_input"

	counterName        metrics.CounterName = "reports"
	counterDescription string              = "the number of reports managed by the reports service"
	topicName          string              = "reports"
	serviceName        string              = "reports_service"
)

var (
	_ models.ReportDataServer = (*Service)(nil)
)

type (
	// Service handles to-do list reports
	Service struct {
		logger            logging.Logger
		reportDataManager models.ReportDataManager
		reportIDFetcher   ReportIDFetcher
		userIDFetcher     UserIDFetcher
		reportCounter     metrics.UnitCounter
		encoderDecoder    encoding.EncoderDecoder
		reporter          newsman.Reporter
	}

	// UserIDFetcher is a function that fetches user IDs.
	UserIDFetcher func(*http.Request) uint64

	// ReportIDFetcher is a function that fetches report IDs.
	ReportIDFetcher func(*http.Request) uint64
)

// ProvideReportsService builds a new ReportsService.
func ProvideReportsService(
	logger logging.Logger,
	reportDataManager models.ReportDataManager,
	reportIDFetcher ReportIDFetcher,
	userIDFetcher UserIDFetcher,
	encoder encoding.EncoderDecoder,
	reportCounterProvider metrics.UnitCounterProvider,
	reporter newsman.Reporter,
) (*Service, error) {
	reportCounter, err := reportCounterProvider(counterName, counterDescription)
	if err != nil {
		return nil, fmt.Errorf("error initializing counter: %w", err)
	}

	svc := &Service{
		logger:            logger.WithName(serviceName),
		reportIDFetcher:   reportIDFetcher,
		userIDFetcher:     userIDFetcher,
		reportDataManager: reportDataManager,
		encoderDecoder:    encoder,
		reportCounter:     reportCounter,
		reporter:          reporter,
	}

	return svc, nil
}
