/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidIngredientMeasurementUnit } from '../models/ValidIngredientMeasurementUnit';
import type { ValidIngredientMeasurementUnitCreationRequestInput } from '../models/ValidIngredientMeasurementUnitCreationRequestInput';
import type { ValidIngredientMeasurementUnitUpdateRequestInput } from '../models/ValidIngredientMeasurementUnitUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidIngredientMeasurementUnitsService {
  /**
   * Operation for fetching ValidIngredientMeasurementUnit
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
  public static getValidIngredientMeasurementUnits(
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
      data?: Array<ValidIngredientMeasurementUnit>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_measurement_units',
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
   * Operation for creating ValidIngredientMeasurementUnit
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidIngredientMeasurementUnit(
    requestBody: ValidIngredientMeasurementUnitCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientMeasurementUnit;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_ingredient_measurement_units',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidIngredientMeasurementUnit
   * @param validIngredientMeasurementUnitId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidIngredientMeasurementUnit(validIngredientMeasurementUnitId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientMeasurementUnit;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}',
      path: {
        validIngredientMeasurementUnitID: validIngredientMeasurementUnitId,
      },
    });
  }
  /**
   * Operation for fetching ValidIngredientMeasurementUnit
   * @param validIngredientMeasurementUnitId
   * @returns any
   * @throws ApiError
   */
  public static getValidIngredientMeasurementUnit(validIngredientMeasurementUnitId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientMeasurementUnit;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}',
      path: {
        validIngredientMeasurementUnitID: validIngredientMeasurementUnitId,
      },
    });
  }
  /**
   * Operation for updating ValidIngredientMeasurementUnit
   * @param validIngredientMeasurementUnitId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidIngredientMeasurementUnit(
    validIngredientMeasurementUnitId: string,
    requestBody: ValidIngredientMeasurementUnitUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientMeasurementUnit;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}',
      path: {
        validIngredientMeasurementUnitID: validIngredientMeasurementUnitId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidIngredientMeasurementUnit
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
  public static getValidIngredientMeasurementUnitsByIngredient(
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
      data?: Array<ValidIngredientMeasurementUnit>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_measurement_units/by_ingredient/{validIngredientID}',
      path: {
        validIngredientID: validIngredientId,
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
   * Operation for fetching ValidIngredientMeasurementUnit
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param validMeasurementUnitId
   * @returns any
   * @throws ApiError
   */
  public static getValidIngredientMeasurementUnitsByMeasurementUnit(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    validMeasurementUnitId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<ValidIngredientMeasurementUnit>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_measurement_units/by_measurement_unit/{validMeasurementUnitID}',
      path: {
        validMeasurementUnitID: validMeasurementUnitId,
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
