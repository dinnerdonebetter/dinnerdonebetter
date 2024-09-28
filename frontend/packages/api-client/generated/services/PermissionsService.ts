/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { UserPermissionsRequestInput } from '../models/UserPermissionsRequestInput';
import type { UserPermissionsResponse } from '../models/UserPermissionsResponse';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class PermissionsService {
  /**
   * Operation for creating UserPermissionsResponse
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static checkPermissions(requestBody: UserPermissionsRequestInput): CancelablePromise<
    APIResponse & {
      data?: UserPermissionsResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/users/permissions/check',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
