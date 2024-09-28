/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidIngredientPreparation } from '../models/ValidIngredientPreparation';
import type { ValidIngredientPreparationCreationRequestInput } from '../models/ValidIngredientPreparationCreationRequestInput';
import type { ValidIngredientPreparationUpdateRequestInput } from '../models/ValidIngredientPreparationUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidIngredientPreparationsService {
  /**
   * Operation for fetching ValidIngredientPreparation
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
  public static getValidIngredientPreparations(
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
      data?: Array<ValidIngredientPreparation>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_preparations',
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
   * Operation for creating ValidIngredientPreparation
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidIngredientPreparation(
    requestBody: ValidIngredientPreparationCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientPreparation;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_ingredient_preparations',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidIngredientPreparation
   * @param validIngredientPreparationId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidIngredientPreparation(validIngredientPreparationId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientPreparation;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_ingredient_preparations/{validIngredientPreparationID}',
      path: {
        validIngredientPreparationID: validIngredientPreparationId,
      },
    });
  }
  /**
   * Operation for fetching ValidIngredientPreparation
   * @param validIngredientPreparationId
   * @returns any
   * @throws ApiError
   */
  public static getValidIngredientPreparation(validIngredientPreparationId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientPreparation;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_preparations/{validIngredientPreparationID}',
      path: {
        validIngredientPreparationID: validIngredientPreparationId,
      },
    });
  }
  /**
   * Operation for updating ValidIngredientPreparation
   * @param validIngredientPreparationId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidIngredientPreparation(
    validIngredientPreparationId: string,
    requestBody: ValidIngredientPreparationUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientPreparation;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_ingredient_preparations/{validIngredientPreparationID}',
      path: {
        validIngredientPreparationID: validIngredientPreparationId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidIngredientPreparation
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
  public static getValidIngredientPreparationsByIngredient(
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
      data?: Array<ValidIngredientPreparation>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_preparations/by_ingredient/{validIngredientID}',
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
   * Operation for fetching ValidIngredientPreparation
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
  public static getValidIngredientPreparationsByPreparation(
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
      data?: Array<ValidIngredientPreparation>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_preparations/by_preparation/{validPreparationID}',
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
