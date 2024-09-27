/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { HouseholdInvitation } from '../models/HouseholdInvitation';
import type { HouseholdInvitationCreationRequestInput } from '../models/HouseholdInvitationCreationRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class InvitationsService {
  /**
   * Operation for creating HouseholdInvitation
   * @param householdId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createHouseholdInvitation(
    householdId: string,
    requestBody: HouseholdInvitationCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: HouseholdInvitation;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/households/{householdID}/invitations',
      path: {
        householdID: householdId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching HouseholdInvitation
   * @param householdId
   * @param householdInvitationId
   * @returns any
   * @throws ApiError
   */
  public static getHouseholdInvitationById(
    householdId: string,
    householdInvitationId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: HouseholdInvitation;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/households/{householdID}/invitations/{householdInvitationID}',
      path: {
        householdID: householdId,
        householdInvitationID: householdInvitationId,
      },
    });
  }
}
