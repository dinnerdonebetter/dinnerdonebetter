/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { TOTPSecretRefreshInput } from '../models/TOTPSecretRefreshInput';
import type { TOTPSecretRefreshResponse } from '../models/TOTPSecretRefreshResponse';
import type { TOTPSecretVerificationInput } from '../models/TOTPSecretVerificationInput';
import type { User } from '../models/User';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class TotpSecretService {
  /**
   * Operation for creating TOTPSecretRefreshResponse
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static refreshTotpSecret(requestBody: TOTPSecretRefreshInput): CancelablePromise<
    APIResponse & {
      data?: TOTPSecretRefreshResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/users/totp_secret/new',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating User
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static verifyTotpSecret(requestBody: TOTPSecretVerificationInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/totp_secret/verify',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
