/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { UserAccountStatusUpdateInput } from '../models/UserAccountStatusUpdateInput';
import type { UserStatusResponse } from '../models/UserStatusResponse';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class AdminService {
  /**
   * Operation for creating
   * @throws ApiError
   */
  public static adminCycleCookieSecret(): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/admin/cycle_cookie_secret',
    });
  }
  /**
   * Operation for creating UserStatusResponse
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static adminUpdateUserStatus(requestBody: UserAccountStatusUpdateInput): CancelablePromise<
    APIResponse & {
      data?: UserStatusResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/admin/users/status',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
