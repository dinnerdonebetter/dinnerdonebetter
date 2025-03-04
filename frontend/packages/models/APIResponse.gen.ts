// GENERATED CODE, DO NOT EDIT MANUALLY

import { IAPIError } from './APIError.gen';
import { ResponseDetails } from './ResponseDetails.gen';
import { Pagination } from './Pagination.gen';

export class APIResponse<T> {
  data: T;
  pagination?: Pagination;
  error?: IAPIError;
  details: ResponseDetails;

  constructor(
    input: {
      data?: T;
      pagination?: Pagination;
      error?: IAPIError;
      details: ResponseDetails;
    } = {
      details: new ResponseDetails(),
    },
  ) {
    this.data = input.data || ({} as T);
    this.pagination = input.pagination;
    this.error = input.error;
    this.details = input.details;
  }
}
