/**
 * Server-side OpenTelemetry tracing for SvelteKit (Node).
 * Import this module from hooks.server.ts so it runs at startup.
 * When OTEL_COLLECTOR_GRPC_URL is unset, tracing is no-op.
 */

import { NodeSDK } from '@opentelemetry/sdk-node';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { Resource } from '@opentelemetry/resources';
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';
import { env } from '$env/dynamic/private';

const SERVICE_NAME = 'consumer-web';

let initialized = false;

export function initServerOtel(): void {
  if (initialized) return;
  initialized = true;

  const collectorGrpcUrl = env.OTEL_COLLECTOR_GRPC_URL?.trim();
  if (!collectorGrpcUrl) {
    return;
  }

  const endpoint =
    collectorGrpcUrl.startsWith('http://') || collectorGrpcUrl.startsWith('https://')
      ? collectorGrpcUrl
      : `http://${collectorGrpcUrl}`;

  const traceExporter = new OTLPTraceExporter({ url: endpoint });
  const resource = new Resource({ [ATTR_SERVICE_NAME]: SERVICE_NAME });

  const sdk = new NodeSDK({
    resource,
    traceExporter,
    // Rely on collector for sampling in prod; trace everything when collector is configured for now.
    sampler: undefined,
  });

  sdk.start();
}
