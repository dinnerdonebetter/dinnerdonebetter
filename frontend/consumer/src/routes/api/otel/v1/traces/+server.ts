import { env } from '$env/dynamic/private';
import type { RequestHandler } from './$types';

/**
 * OTLP proxy for browser trace export. The browser sends OTLP HTTP (POST) to this
 * same-origin route; we forward to the collector's OTLP HTTP endpoint so the
 * collector URL stays server-only and CORS is avoided.
 */
export const POST: RequestHandler = async ({ request }) => {
  const collectorHttpUrl = env.OTEL_COLLECTOR_HTTP_URL?.trim();
  if (!collectorHttpUrl) {
    return new Response(null, { status: 501 });
  }

  const base = collectorHttpUrl.replace(/\/$/, '');
  const tracesUrl = `${base}/v1/traces`;

  const body = await request.arrayBuffer();
  const contentType = request.headers.get('content-type') ?? 'application/x-protobuf';

  const res = await fetch(tracesUrl, {
    method: 'POST',
    headers: {
      'Content-Type': contentType,
    },
    body: body.byteLength > 0 ? body : undefined,
  });

  return new Response(res.body, {
    status: res.status,
    headers: {
      'Content-Type': res.headers.get('content-type') ?? 'application/json',
    },
  });
};
