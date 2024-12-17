package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics/otelgrpc"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func setupTracingProvider(ctx context.Context, collectorAddress string) (trace.TracerProvider, func(), error) {
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithOSType(),
		resource.WithAttributes(
			attribute.KeyValue{
				Key:   "service.name",
				Value: attribute.StringValue("demo-server"),
			},
		),
	)
	if err != nil {
		return nil, nil, err
	}

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(collectorAddress),
	)
	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	shutdownFunc := func() {
		cxt, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err = traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}

	return tracerProvider, shutdownFunc, nil
}

func setupMetricsProvider(ctx context.Context, collectorAddress string) (metric.MeterProvider, func(), error) {
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithOSType(),
		resource.WithAttributes(
			attribute.KeyValue{
				Key:   "service.name",
				Value: attribute.StringValue("demo-server"),
			},
		),
	)
	if err != nil {
		return nil, nil, err
	}

	metricExp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(collectorAddress),
	)
	handleErr(err, "Failed to create the collector metric exporter")

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOnFilter),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExp,
				sdkmetric.WithInterval(2*time.Second),
			),
		),
	)
	otel.SetMeterProvider(meterProvider)

	shutdownFunc := func() {
		cxt, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err = meterProvider.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}

	return meterProvider, shutdownFunc, nil
}

func initProvider() (metric.MeterProvider, trace.TracerProvider, func()) {
	ctx := context.Background()

	otelAgentAddr, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if !ok {
		otelAgentAddr = "otel_collector:4317"
	}

	meterProvider, metricShutdown, err := setupMetricsProvider(ctx, otelAgentAddr)
	handleErr(err, "Failed to create the metrics provider")

	tracerProvider, tracerShutdown, err := setupTracingProvider(ctx, otelAgentAddr)
	handleErr(err, "Failed to create the tracer provider")

	return meterProvider, tracerProvider, func() {
		tracerShutdown()
		metricShutdown()
	}
}

func initNewProvider(logger logging.Logger) (metrics.Provider, trace.TracerProvider, func()) {
	ctx := context.Background()

	otelAgentAddr, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if !ok {
		otelAgentAddr = "otel_collector:4317"
	}

	meterProviderZ, err := otelgrpc.ProvideMetricsProvider(ctx, logger, &otelgrpc.Config{
		BaseName:           "api-service",
		CollectorEndpoint:  otelAgentAddr,
		CollectionInterval: time.Second,
		Insecure:           true,
		CollectionTimeout:  5 * time.Second,
	})
	handleErr(err, "Failed to create the metrics provider")

	tracerProvider, tracerShutdown, err := setupTracingProvider(ctx, otelAgentAddr)
	handleErr(err, "Failed to create the tracer provider")

	return meterProviderZ, tracerProvider, func() {
		tracerShutdown()
		meterProviderZ.Shutdown(context.Background())
	}
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

type errHandler struct{}

func (e *errHandler) Handle(err error) {
	log.Println(err)
}

func buildCounters(useNew bool, logger logging.Logger) (metric.Int64Counter, metric.Int64Counter, func()) {
	if useNew {
		meterProvider, _, shutdown := initNewProvider(logger)

		otel.SetErrorHandler(&errHandler{})

		requestCount, err := meterProvider.NewInt64Counter(
			"demo_server/request_counts",
			metric.WithDescription("The number of requests received"),
		)
		handleErr(err, "failed to create request count metric")

		arbitraryCount, err := meterProvider.NewInt64Counter(
			"arbitrary",
			metric.WithDescription("Meaningless number"),
		)
		handleErr(err, "failed to create request count metric")

		return requestCount, arbitraryCount, shutdown
	} else {
		meterProvider, _, shutdown := initProvider()

		otel.SetErrorHandler(&errHandler{})

		meter := meterProvider.Meter("demo-server-meter", metric.WithInstrumentationAttributes(attribute.KeyValue{
			Key:   "service.name",
			Value: attribute.StringValue("demo-server"),
		}))

		requestCount, err := meter.Int64Counter(
			"demo_server/request_counts",
			metric.WithDescription("The number of requests received"),
		)
		handleErr(err, "failed to create request count metric")

		arbitraryCount, err := meter.Int64Counter(
			"arbitrary",
			metric.WithDescription("Meaningless number"),
		)
		handleErr(err, "failed to create request count metric")

		return requestCount, arbitraryCount, shutdown
	}
}

func main() {
	logger := loggingcfg.ProvideLogger(&loggingcfg.Config{
		Provider:       loggingcfg.ProviderSlog,
		OutputFilepath: "/var/log/dinnerdonebetter/demo-server.log",
	})

	serverAttribute := attribute.String("server-attribute", "foo")
	commonLabels := []attribute.KeyValue{serverAttribute}

	otel.SetErrorHandler(&errHandler{})

	requestCount, arbitraryCount, shutdown := buildCounters(false, logger)
	defer shutdown()

	go func() {
		for {
			logger.Info("arbitrary message!")
			arbitraryCount.Add(context.Background(), 1, metric.WithAttributes(commonLabels...))
			time.Sleep(time.Second)
		}
	}()

	// create a handler wrapped in OpenTelemetry instrumentation
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//  random sleep to simulate latency
		var sleep int64

		switch modulus := time.Now().Unix() % 5; modulus {
		case 0:
			sleep = rng.Int63n(2048)
		case 1:
			sleep = rng.Int63n(16)
		case 2:
			sleep = rng.Int63n(512)
		case 3:
			sleep = rng.Int63n(128)
		case 4:
			sleep = rng.Int63n(1024)
		}
		time.Sleep(time.Duration(sleep) * time.Millisecond)

		ctx := req.Context()
		requestCount.Add(ctx, 1, metric.WithAttributes(commonLabels...))
		span := trace.SpanFromContext(ctx)
		bag := baggage.FromContext(ctx)

		baggageAttributes := []attribute.KeyValue{serverAttribute}
		for _, member := range bag.Members() {
			baggageAttributes = append(baggageAttributes, attribute.String("baggage key:"+member.Key(), member.Value()))
		}
		span.SetAttributes(baggageAttributes...)

		if _, err := w.Write([]byte("Hello World\n")); err != nil {
			http.Error(w, "write operation failed.", http.StatusInternalServerError)
			return
		}
	})

	mux := http.NewServeMux()
	mux.Handle("/hello", otelhttp.NewHandler(handler, "/hello"))

	server := &http.Server{
		Addr:              ":8000",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		handleErr(err, "server failed to serve")
	}
}
