/**
 * Client-side OpenTelemetry tracing and metrics for the browser.
 * Spans are sent to the same-origin OTLP proxy at /api/otel/v1/traces.
 * Metrics are sent to /api/otel/v1/metrics.
 * No-op when run on the server (SSR).
 */

import { WebTracerProvider, BatchSpanProcessor, TraceIdRatioBasedSampler } from '@opentelemetry/sdk-trace-web';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { MeterProvider, PeriodicExportingMetricReader } from '@opentelemetry/sdk-metrics';
import { OTLPMetricExporter } from '@opentelemetry/exporter-metrics-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';
import { metrics } from '@opentelemetry/api';
import { shouldRecordPathForMetrics } from './utils';

const SERVICE_NAME = 'consumer-web';
const METER_NAME = 'consumer-web-client';
const METER_VERSION = '1.0';

const EXPORT_INTERVAL_MS = 60_000;

let tracingInitialized = false;
let metricsInitialized = false;

export function initClientOtel(): void {
  if (typeof window === 'undefined' || tracingInitialized) return;
  tracingInitialized = true;

  const exporter = new OTLPTraceExporter({
    url: `${window.location.origin}/api/otel/v1/traces`,
  });

  const resource = new Resource({ [ATTR_SERVICE_NAME]: SERVICE_NAME });
  const sampler = new TraceIdRatioBasedSampler(1.0);

  const provider = new WebTracerProvider({
    resource,
    sampler,
    spanProcessors: [new BatchSpanProcessor(exporter)],
  });

  provider.register();
}

/**
 * Initialize client-side metrics. Sends to same-origin proxy at /api/otel/v1/metrics.
 * Call from +layout.svelte alongside initClientOtel() when in browser.
 */
export function initClientOtelMetrics(): void {
  if (typeof window === 'undefined' || metricsInitialized) return;
  metricsInitialized = true;

  const resource = new Resource({ [ATTR_SERVICE_NAME]: SERVICE_NAME });
  const exporter = new OTLPMetricExporter({
    url: `${window.location.origin}/api/otel/v1/metrics`,
  });

  const reader = new PeriodicExportingMetricReader({
    exporter,
    exportIntervalMillis: EXPORT_INTERVAL_MS,
  });

  const meterProvider = new MeterProvider({
    resource,
    readers: [reader],
  });

  metrics.setGlobalMeterProvider(meterProvider);

  const meter = meterProvider.getMeter(METER_NAME, METER_VERSION);
  _pageViewCounter = meter.createCounter('consumer_web_page_views', {
    description: 'Number of page views',
  });
}

let _pageViewCounter: { add: (value: number, attributes?: Record<string, string>) => void } | null =
  null;

/**
 * Record a page view for metrics. Call from the client after navigation (e.g. afterNavigate).
 * No-op if metrics are not initialized. Single counter (no route label) to keep cardinality low.
 */
export function recordPageView(route: string): void {
  if (typeof window === 'undefined') return;
  if (!shouldRecordPathForMetrics(route)) return;
  _pageViewCounter?.add(1);
}
