/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { Household } from '../models/Household';
import type { HouseholdCreationRequestInput } from '../models/HouseholdCreationRequestInput';
import type { HouseholdInstrumentOwnership } from '../models/HouseholdInstrumentOwnership';
import type { HouseholdInstrumentOwnershipCreationRequestInput } from '../models/HouseholdInstrumentOwnershipCreationRequestInput';
import type { HouseholdInstrumentOwnershipUpdateRequestInput } from '../models/HouseholdInstrumentOwnershipUpdateRequestInput';
import type { HouseholdInvitation } from '../models/HouseholdInvitation';
import type { HouseholdInvitationCreationRequestInput } from '../models/HouseholdInvitationCreationRequestInput';
import type { HouseholdOwnershipTransferInput } from '../models/HouseholdOwnershipTransferInput';
import type { HouseholdUpdateRequestInput } from '../models/HouseholdUpdateRequestInput';
import type { HouseholdUserMembership } from '../models/HouseholdUserMembership';
import type { ModifyUserPermissionsInput } from '../models/ModifyUserPermissionsInput';
import type { UserPermissionsResponse } from '../models/UserPermissionsResponse';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class HouseholdsService {
  /**
   * Operation for fetching Household
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
  public static getHouseholds(
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
      data?: Array<Household>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/households',
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
   * Operation for creating Household
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createHousehold(requestBody: HouseholdCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: Household;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/households',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving Household
   * @param householdId
   * @returns any
   * @throws ApiError
   */
  public static archiveHousehold(householdId: string): CancelablePromise<
    APIResponse & {
      data?: Household;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/households/{householdID}',
      path: {
        householdID: householdId,
      },
    });
  }
  /**
   * Operation for fetching Household
   * @param householdId
   * @returns any
   * @throws ApiError
   */
  public static getHousehold(householdId: string): CancelablePromise<
    APIResponse & {
      data?: Household;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/households/{householdID}',
      path: {
        householdID: householdId,
      },
    });
  }
  /**
   * Operation for updating Household
   * @param householdId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateHousehold(
    householdId: string,
    requestBody: HouseholdUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: Household;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/households/{householdID}',
      path: {
        householdID: householdId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating Household
   * @param householdId
   * @returns any
   * @throws ApiError
   */
  public static setDefaultHousehold(householdId: string): CancelablePromise<
    APIResponse & {
      data?: Household;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/households/{householdID}/default',
      path: {
        householdID: householdId,
      },
    });
  }
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
  /**
   * Operation for creating HouseholdInvitation
   * @param householdId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static postHouseholdsHouseholdIdInvite(
    householdId: string,
    requestBody: HouseholdInvitationCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: HouseholdInvitation;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/households/{householdID}/invite',
      path: {
        householdID: householdId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
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
  /**
   * Operation for creating Household
   * @param householdId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static transferHouseholdOwnership(
    householdId: string,
    requestBody: HouseholdOwnershipTransferInput,
  ): CancelablePromise<
    APIResponse & {
      data?: Household;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/households/{householdID}/transfer',
      path: {
        householdID: householdId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching Household
   * @returns any
   * @throws ApiError
   */
  public static getActiveHousehold(): CancelablePromise<
    APIResponse & {
      data?: Household;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/households/current',
    });
  }
  /**
   * Operation for fetching HouseholdInstrumentOwnership
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
  public static getHouseholdInstrumentOwnerships(
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
      data?: Array<HouseholdInstrumentOwnership>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/households/instruments',
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
   * Operation for creating HouseholdInstrumentOwnership
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createHouseholdInstrumentOwnership(
    requestBody: HouseholdInstrumentOwnershipCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: HouseholdInstrumentOwnership;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/households/instruments',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving HouseholdInstrumentOwnership
   * @param householdInstrumentOwnershipId
   * @returns any
   * @throws ApiError
   */
  public static archiveHouseholdInstrumentOwnership(householdInstrumentOwnershipId: string): CancelablePromise<
    APIResponse & {
      data?: HouseholdInstrumentOwnership;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/households/instruments/{householdInstrumentOwnershipID}',
      path: {
        householdInstrumentOwnershipID: householdInstrumentOwnershipId,
      },
    });
  }
  /**
   * Operation for fetching HouseholdInstrumentOwnership
   * @param householdInstrumentOwnershipId
   * @returns any
   * @throws ApiError
   */
  public static getHouseholdInstrumentOwnership(householdInstrumentOwnershipId: string): CancelablePromise<
    APIResponse & {
      data?: HouseholdInstrumentOwnership;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/households/instruments/{householdInstrumentOwnershipID}',
      path: {
        householdInstrumentOwnershipID: householdInstrumentOwnershipId,
      },
    });
  }
  /**
   * Operation for updating HouseholdInstrumentOwnership
   * @param householdInstrumentOwnershipId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateHouseholdInstrumentOwnership(
    householdInstrumentOwnershipId: string,
    requestBody: HouseholdInstrumentOwnershipUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: HouseholdInstrumentOwnership;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/households/instruments/{householdInstrumentOwnershipID}',
      path: {
        householdInstrumentOwnershipID: householdInstrumentOwnershipId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
