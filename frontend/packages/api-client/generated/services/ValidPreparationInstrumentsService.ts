/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidPreparationInstrument } from '../models/ValidPreparationInstrument';
import type { ValidPreparationInstrumentCreationRequestInput } from '../models/ValidPreparationInstrumentCreationRequestInput';
import type { ValidPreparationInstrumentUpdateRequestInput } from '../models/ValidPreparationInstrumentUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidPreparationInstrumentsService {
  /**
   * Operation for fetching ValidPreparationInstrument
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
  public static getValidPreparationInstruments(
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
      data?: Array<ValidPreparationInstrument>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparation_instruments',
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
   * Operation for creating ValidPreparationInstrument
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidPreparationInstrument(
    requestBody: ValidPreparationInstrumentCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidPreparationInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_preparation_instruments',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidPreparationInstrument
   * @param validPreparationVesselId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidPreparationInstrument(validPreparationVesselId: string): CancelablePromise<
    APIResponse & {
      data?: ValidPreparationInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_preparation_instruments/{validPreparationVesselID}',
      path: {
        validPreparationVesselID: validPreparationVesselId,
      },
    });
  }
  /**
   * Operation for fetching ValidPreparationInstrument
   * @param validPreparationVesselId
   * @returns any
   * @throws ApiError
   */
  public static getValidPreparationInstrument(validPreparationVesselId: string): CancelablePromise<
    APIResponse & {
      data?: ValidPreparationInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparation_instruments/{validPreparationVesselID}',
      path: {
        validPreparationVesselID: validPreparationVesselId,
      },
    });
  }
  /**
   * Operation for updating ValidPreparationInstrument
   * @param validPreparationVesselId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidPreparationInstrument(
    validPreparationVesselId: string,
    requestBody: ValidPreparationInstrumentUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidPreparationInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_preparation_instruments/{validPreparationVesselID}',
      path: {
        validPreparationVesselID: validPreparationVesselId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidPreparationInstrument
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param validInstrumentId
   * @returns any
   * @throws ApiError
   */
  public static getValidPreparationInstrumentsByInstrument(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    validInstrumentId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<ValidPreparationInstrument>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparation_instruments/by_instrument/{validInstrumentID}',
      path: {
        validInstrumentID: validInstrumentId,
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
   * Operation for fetching ValidPreparationInstrument
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
  public static getValidPreparationInstrumentsByPreparation(
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
      data?: Array<ValidPreparationInstrument>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_preparation_instruments/by_preparation/{validPreparationID}',
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
}
