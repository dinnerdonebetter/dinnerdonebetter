/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidPreparationVessel } from '../models/ValidPreparationVessel';
import type { ValidPreparationVesselCreationRequestInput } from '../models/ValidPreparationVesselCreationRequestInput';
import type { ValidPreparationVesselUpdateRequestInput } from '../models/ValidPreparationVesselUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidPreparationVesselsService {
  /**
   * Operation for fetching ValidPreparationVessel
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
  public static getValidPreparationVessels(
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
      data?: Array<ValidPreparationVessel>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparation_vessels',
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
   * Operation for creating ValidPreparationVessel
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidPreparationVessel(
    requestBody: ValidPreparationVesselCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidPreparationVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_preparation_vessels',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidPreparationVessel
   * @param validPreparationVesselId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidPreparationVessel(validPreparationVesselId: string): CancelablePromise<
    APIResponse & {
      data?: ValidPreparationVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_preparation_vessels/{validPreparationVesselID}',
      path: {
        validPreparationVesselID: validPreparationVesselId,
      },
    });
  }
  /**
   * Operation for fetching ValidPreparationVessel
   * @param validPreparationVesselId
   * @returns any
   * @throws ApiError
   */
  public static getValidPreparationVessel(validPreparationVesselId: string): CancelablePromise<
    APIResponse & {
      data?: ValidPreparationVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparation_vessels/{validPreparationVesselID}',
      path: {
        validPreparationVesselID: validPreparationVesselId,
      },
    });
  }
  /**
   * Operation for updating ValidPreparationVessel
   * @param validPreparationVesselId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidPreparationVessel(
    validPreparationVesselId: string,
    requestBody: ValidPreparationVesselUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidPreparationVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_preparation_vessels/{validPreparationVesselID}',
      path: {
        validPreparationVesselID: validPreparationVesselId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidPreparationVessel
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param validPreparationId
   * @returns any
   * @throws ApiError
   */
  public static getValidPreparationVesselsByPreparation(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    validPreparationId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<ValidPreparationVessel>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparation_vessels/by_preparation/{validPreparationID}',
      path: {
        validPreparationID: validPreparationId,
      },
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
   * Operation for fetching ValidPreparationVessel
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param validVesselId
   * @returns any
   * @throws ApiError
   */
  public static getValidPreparationVesselsByVessel(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    validVesselId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<ValidPreparationVessel>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparation_vessels/by_vessel/{ValidVesselID}',
      path: {
        ValidVesselID: validVesselId,
      },
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
