/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIError } from './APIError';
import type { Pagination } from './Pagination';
import type { ResponseDetails } from './ResponseDetails';
export type APIResponse = {
  details?: ResponseDetails;
  error?: APIError;
  pagination?: Pagination;
};
