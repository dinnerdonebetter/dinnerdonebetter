/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { UserStatusResponse } from '../models/UserStatusResponse';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class AuthService {
  /**
   * Operation for fetching UserStatusResponse
   * @returns any
   * @throws ApiError
   */
  public static getAuthStatus(): CancelablePromise<
    APIResponse & {
      data?: UserStatusResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/auth/status',
    });
  }
}
