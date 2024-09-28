/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { HouseholdInvitation } from '../models/HouseholdInvitation';
import type { HouseholdInvitationUpdateRequestInput } from '../models/HouseholdInvitationUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class HouseholdInvitationsService {
  /**
   * Operation for fetching HouseholdInvitation
   * @param householdInvitationId
   * @returns any
   * @throws ApiError
   */
  public static getHouseholdInvitation(householdInvitationId: string): CancelablePromise<
    APIResponse & {
      data?: HouseholdInvitation;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/household_invitations/{householdInvitationID}',
      path: {
        householdInvitationID: householdInvitationId,
      },
    });
  }
  /**
   * Operation for updating HouseholdInvitation
   * @param householdInvitationId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static acceptHouseholdInvitation(
    householdInvitationId: string,
    requestBody: HouseholdInvitationUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: HouseholdInvitation;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/household_invitations/{householdInvitationID}/accept',
      path: {
        householdInvitationID: householdInvitationId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for updating HouseholdInvitation
   * @param householdInvitationId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static cancelHouseholdInvitation(
    householdInvitationId: string,
    requestBody: HouseholdInvitationUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: HouseholdInvitation;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/household_invitations/{householdInvitationID}/cancel',
      path: {
        householdInvitationID: householdInvitationId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for updating HouseholdInvitation
   * @param householdInvitationId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static rejectHouseholdInvitation(
    householdInvitationId: string,
    requestBody: HouseholdInvitationUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: HouseholdInvitation;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/household_invitations/{householdInvitationID}/reject',
      path: {
        householdInvitationID: householdInvitationId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching HouseholdInvitation
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @returns any
   * @throws ApiError
   */
  public static getReceivedHouseholdInvitations(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
  ): CancelablePromise<
    APIResponse & {
      data?: Array<HouseholdInvitation>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/household_invitations/received',
      query: {
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
  /**
   * Operation for fetching HouseholdInvitation
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @returns any
   * @throws ApiError
   */
  public static getSentHouseholdInvitations(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
  ): CancelablePromise<
    APIResponse & {
      data?: Array<HouseholdInvitation>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/household_invitations/sent',
      query: {
        limit: limit,
        page: page,
        createdBefore: createdBefore,
        createdAfter: createdAfter,
        updatedBefore: updatedBefore,
        updatedAfter: updatedAfter,
        includeArchived: includeArchived,
        sortBy: sortBy,
      },
    });
  }
}
