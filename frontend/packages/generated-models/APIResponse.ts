// GENERATED CODE, DO NOT EDIT MANUALLY

import { APIError } from './APIError';
import { Pagination } from './Pagination';
import { ResponseDetails } from './ResponseDetails';

export interface IAPIResponse {
  pagination?: Pagination;
  details: ResponseDetails;
  error?: APIError;
}

export class APIResponse implements IAPIResponse {
  pagination?: Pagination;
  details: ResponseDetails;
  error?: APIError;
  constructor(input: Partial<APIResponse> = {}) {
    this.pagination = input.pagination;
    this.details = input.details = new ResponseDetails();
    this.error = input.error;
  }
}
