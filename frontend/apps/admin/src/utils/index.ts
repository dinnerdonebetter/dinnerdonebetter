import { EitherErrorOr } from '@dinnerdonebetter/models';

export function valueOrDefault<T>(x: EitherErrorOr<T>, defaultValue: any): T {
  return x.error ? defaultValue : x.data ? x.data : defaultValue;
}
