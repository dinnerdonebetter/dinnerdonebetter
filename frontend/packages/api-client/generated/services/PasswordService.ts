/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { PasswordResetToken } from '../models/PasswordResetToken';
import type { PasswordResetTokenCreationRequestInput } from '../models/PasswordResetTokenCreationRequestInput';
import type { PasswordResetTokenRedemptionRequestInput } from '../models/PasswordResetTokenRedemptionRequestInput';
import type { PasswordUpdateInput } from '../models/PasswordUpdateInput';
import type { User } from '../models/User';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class PasswordService {
  /**
   * Operation for updating
   * @param requestBody
   * @throws ApiError
   */
  public static updatePassword(requestBody: PasswordUpdateInput): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/users/password/new',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating PasswordResetToken
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static requestPasswordResetToken(requestBody: PasswordResetTokenCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: PasswordResetToken;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/password/reset',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for redeeming a password reset token
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static redeemPasswordResetToken(requestBody: PasswordResetTokenRedemptionRequestInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/password/reset/redeem',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
