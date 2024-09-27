/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { EmailAddressVerificationRequestInput } from '../models/EmailAddressVerificationRequestInput';
import type { User } from '../models/User';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class EmailAddressService {
  /**
   * Operation for creating User
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static verifyEmailAddress(requestBody: EmailAddressVerificationRequestInput): CancelablePromise<
    APIResponse & {
      data?: User;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/users/email_address/verify',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
