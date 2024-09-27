/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { HouseholdUserMembership } from '../models/HouseholdUserMembership';
import type { ModifyUserPermissionsInput } from '../models/ModifyUserPermissionsInput';
import type { UserPermissionsResponse } from '../models/UserPermissionsResponse';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class HouseholdMembersService {
  /**
   * Operation for archiving HouseholdUserMembership
   * @param householdId
   * @param userId
   * @returns any
   * @throws ApiError
   */
  public static archiveUserMembership(
    householdId: string,
    userId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: HouseholdUserMembership;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/households/{householdID}/members/{userID}',
      path: {
        householdID: householdId,
        userID: userId,
      },
    });
  }
  /**
   * Operation for updating UserPermissionsResponse
   * @param householdId
   * @param userId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateHouseholdMemberPermissions(
    householdId: string,
    userId: string,
    requestBody: ModifyUserPermissionsInput,
  ): CancelablePromise<
    APIResponse & {
      data?: UserPermissionsResponse;
    }
  > {
    return __request(OpenAPI, {
      method: 'PATCH',
      url: '/api/v1/households/{householdID}/members/{userID}/permissions',
      path: {
        householdID: householdId,
        userID: userId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
