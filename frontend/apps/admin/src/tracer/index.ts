import { Tracer } from '@opentelemetry/api';
import { getTracer } from '@dinnerdonebetter/tracing';

export const serverSideTracer: Tracer = getTracer('admin-server');
