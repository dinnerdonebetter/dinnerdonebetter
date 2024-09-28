/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidVessel } from '../models/ValidVessel';
import type { ValidVesselCreationRequestInput } from '../models/ValidVesselCreationRequestInput';
import type { ValidVesselUpdateRequestInput } from '../models/ValidVesselUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidVesselsService {
  /**
   * Operation for fetching ValidVessel
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
  public static getValidVessels(
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
      data?: Array<ValidVessel>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_vessels',
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
   * Operation for creating ValidVessel
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidVessel(requestBody: ValidVesselCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: ValidVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_vessels',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidVessel
   * @param validVesselId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidVessel(validVesselId: string): CancelablePromise<
    APIResponse & {
      data?: ValidVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_vessels/{validVesselID}',
      path: {
        validVesselID: validVesselId,
      },
    });
  }
  /**
   * Operation for fetching ValidVessel
   * @param validVesselId
   * @returns any
   * @throws ApiError
   */
  public static getValidVessel(validVesselId: string): CancelablePromise<
    APIResponse & {
      data?: ValidVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_vessels/{validVesselID}',
      path: {
        validVesselID: validVesselId,
      },
    });
  }
  /**
   * Operation for updating ValidVessel
   * @param validVesselId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidVessel(
    validVesselId: string,
    requestBody: ValidVesselUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_vessels/{validVesselID}',
      path: {
        validVesselID: validVesselId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidVessel
   * @returns any
   * @throws ApiError
   */
  public static getRandomValidVessel(): CancelablePromise<
    APIResponse & {
      data?: ValidVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_vessels/random',
    });
  }
  /**
   * Operation for fetching ValidVessel
   * @param q the search query parameter
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
  public static searchForValidVessels(
    q: string,
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
      data?: Array<ValidVessel>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_vessels/search',
      query: {
        q: q,
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
