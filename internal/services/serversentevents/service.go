package serversentevents

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/r3labs/sse"
)

const (
	serviceName string = "workers_service"
)

var _ types.ServerSentEventsService = (*service)(nil)

type (
	// service handles valid vessels.
	service struct {
		cfg                       *Config
		logger                    logging.Logger
		dataManager               database.DataManager
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
		stopChan                  chan bool
		eventsServer              *sse.Server
	}
)

// ProvideService builds a new ValidVesselsService.
func ProvideService(
	ctx context.Context,
	cfg *Config,
	logger logging.Logger,
	dataManager database.DataManager,
	encoder encoding.ServerEncoderDecoder,
	consumerProvider messagequeue.ConsumerProvider,
	tracerProvider tracing.TracerProvider,
) (types.ServerSentEventsService, error) {
	sseServer := sse.New()
	sseServer.AutoStream = true

	dataChangesConsumer, err := consumerProvider.ProvideConsumer(ctx, cfg.DataChangesTopicName, func(consumerContext context.Context, payload []byte) error {
		var dataChange types.DataChangeMessage
		if err := encoder.DecodeBytes(consumerContext, payload, &dataChange); err != nil {
			return err
		}

		if dataChange.UserID == "" && sseServer.StreamExists(dataChange.UserID) {
			sseServer.Publish(dataChange.UserID, &sse.Event{
				Data: payload,
			})
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("setting up server-sent events service data changes consumer: %w", err)
	}

	svc := &service{
		cfg:                       cfg,
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		encoderDecoder:            encoder,
		dataManager:               dataManager,
		eventsServer:              sseServer,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	errors := make(chan error)
	go dataChangesConsumer.Consume(svc.stopChan, errors)
	go func() {
		for receivedErr := range errors {
			logger.WithValue("topic_name", cfg.DataChangesTopicName).Error(receivedErr, "consuming data changes")
		}
	}()

	return svc, nil
}
