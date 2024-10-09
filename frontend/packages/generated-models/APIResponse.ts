// GENERATED CODE, DO NOT EDIT MANUALLY

import { APIError } from './APIError';
import { Pagination } from './Pagination';
import { ResponseDetails } from './ResponseDetails';

export interface IAPIResponse {
  details: ResponseDetails;
  error?: APIError;
  pagination?: Pagination;
}

export class APIResponse implements IAPIResponse {
  details: ResponseDetails;
  error?: APIError;
  pagination?: Pagination;
  constructor(input: Partial<APIResponse> = {}) {
    this.details = input.details = new ResponseDetails();
    this.error = input.error;
    this.pagination = input.pagination;
  }
}
