import opentelemetry, { Tracer } from '@opentelemetry/api';
import { AlwaysOnSampler, NodeTracerProvider } from '@opentelemetry/sdk-trace-node';
import { SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { TraceExporter } from '@google-cloud/opentelemetry-cloud-trace-exporter';

// Enable OpenTelemetry exporters to export traces to Google Cloud Trace.
// Exporters use Application Default Credentials (ADCs) to authenticate.
// See https://developers.google.com/identity/protocols/application-default-credentials
// for more details.
const provider = new NodeTracerProvider({
  sampler: new AlwaysOnSampler(),
});

// Initialize the exporter. When your application is running on Google Cloud,
// you don't need to provide auth credentials or a project id.
const exporter = new TraceExporter();

// Configure the span processor to send spans to the exporter
provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

opentelemetry.trace.setGlobalTracerProvider(provider);

export function getTracer(name: string): Tracer {
  return opentelemetry.trace.getTracer(name);
}

export type TracerType = Tracer;
