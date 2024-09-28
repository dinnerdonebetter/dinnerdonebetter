/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ChangeActiveHouseholdInput } from '../models/ChangeActiveHouseholdInput';
import type { Household } from '../models/Household';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class HouseholdService {
  /**
   * Operation for creating Household
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static changeActiveHousehold(requestBody: ChangeActiveHouseholdInput): CancelablePromise<
    APIResponse & {
      data?: Household;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/users/household/select',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
