import { EitherErrorOr } from '@dinnerdonebetter/models';

export function errorOrDefault<T>(x: EitherErrorOr<T>, defaultValue: any): T {
  return x.error ? defaultValue : x.data ? x.data : defaultValue;
}
