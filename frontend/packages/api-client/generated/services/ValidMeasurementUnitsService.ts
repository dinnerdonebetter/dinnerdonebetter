/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidMeasurementUnit } from '../models/ValidMeasurementUnit';
import type { ValidMeasurementUnitCreationRequestInput } from '../models/ValidMeasurementUnitCreationRequestInput';
import type { ValidMeasurementUnitUpdateRequestInput } from '../models/ValidMeasurementUnitUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidMeasurementUnitsService {
  /**
   * Operation for fetching ValidMeasurementUnit
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
  public static getValidMeasurementUnits(
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
      data?: Array<ValidMeasurementUnit>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_measurement_units',
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
   * Operation for creating ValidMeasurementUnit
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidMeasurementUnit(requestBody: ValidMeasurementUnitCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: ValidMeasurementUnit;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_measurement_units',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidMeasurementUnit
   * @param validMeasurementUnitId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidMeasurementUnit(validMeasurementUnitId: string): CancelablePromise<
    APIResponse & {
      data?: ValidMeasurementUnit;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_measurement_units/{validMeasurementUnitID}',
      path: {
        validMeasurementUnitID: validMeasurementUnitId,
      },
    });
  }
  /**
   * Operation for fetching ValidMeasurementUnit
   * @param validMeasurementUnitId
   * @returns any
   * @throws ApiError
   */
  public static getValidMeasurementUnit(validMeasurementUnitId: string): CancelablePromise<
    APIResponse & {
      data?: ValidMeasurementUnit;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_measurement_units/{validMeasurementUnitID}',
      path: {
        validMeasurementUnitID: validMeasurementUnitId,
      },
    });
  }
  /**
   * Operation for updating ValidMeasurementUnit
   * @param validMeasurementUnitId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidMeasurementUnit(
    validMeasurementUnitId: string,
    requestBody: ValidMeasurementUnitUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidMeasurementUnit;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_measurement_units/{validMeasurementUnitID}',
      path: {
        validMeasurementUnitID: validMeasurementUnitId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidMeasurementUnit
   * @param q the search query parameter
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param validIngredientId
   * @returns any
   * @throws ApiError
   */
  public static searchValidMeasurementUnitsByIngredient(
    q: string,
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    validIngredientId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<ValidMeasurementUnit>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_measurement_units/by_ingredient/{validIngredientID}',
      path: {
        validIngredientID: validIngredientId,
      },
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
  /**
   * Operation for fetching ValidMeasurementUnit
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
  public static searchForValidMeasurementUnits(
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
      data?: Array<ValidMeasurementUnit>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_measurement_units/search',
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
