// GENERATED CODE, DO NOT EDIT MANUALLY

import { IAPIError } from './APIError.gen';

export interface EitherErrorOr<T> {
  error?: IAPIError;
  data?: T;
}
