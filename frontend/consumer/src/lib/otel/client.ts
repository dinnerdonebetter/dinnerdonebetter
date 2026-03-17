/**
 * Client-side OpenTelemetry tracing for the browser.
 * Spans are sent to the same-origin OTLP proxy at /api/otel/v1/traces.
 * No-op when run on the server (SSR).
 */

import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-web';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { Resource } from '@opentelemetry/resources';
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';
import { TraceIdRatioBasedSampler } from '@opentelemetry/sdk-trace-web';

const SERVICE_NAME = 'consumer-web';

let initialized = false;

export function initClientOtel(): void {
  if (typeof window === 'undefined' || initialized) return;
  initialized = true;

  const exporter = new OTLPTraceExporter({
    url: '/api/otel/v1/traces',
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
