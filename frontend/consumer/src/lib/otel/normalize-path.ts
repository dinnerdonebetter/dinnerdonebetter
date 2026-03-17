/**
 * Normalize path for metrics to reduce cardinality: replace xids and numeric IDs with :id.
 * Use for path/route attributes on OTel metrics so we don't create a series per resource.
 */

/** xid: 20-char base32hex [0-9a-v] (e.g. github.com/rs/xid) */
const XID_REGEX = /^[0-9a-v]{20}$/;
const NUMERIC_ID_REGEX = /^\d+$/;

export function normalizePathForMetrics(pathname: string): string {
  return pathname
    .split('/')
    .map((segment) => {
      if (!segment) return segment;
      if (XID_REGEX.test(segment)) return ':id';
      if (NUMERIC_ID_REGEX.test(segment)) return ':id';
      return segment;
    })
    .join('/');
}
