/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidIngredient } from '../models/ValidIngredient';
import type { ValidIngredientCreationRequestInput } from '../models/ValidIngredientCreationRequestInput';
import type { ValidIngredientUpdateRequestInput } from '../models/ValidIngredientUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidIngredientsService {
  /**
   * Operation for fetching ValidIngredient
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
  public static getValidIngredients(
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
      data?: Array<ValidIngredient>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredients',
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
   * Operation for creating ValidIngredient
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidIngredient(requestBody: ValidIngredientCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: ValidIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_ingredients',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidIngredient
   * @param validIngredientId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidIngredient(validIngredientId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_ingredients/{validIngredientID}',
      path: {
        validIngredientID: validIngredientId,
      },
    });
  }
  /**
   * Operation for fetching ValidIngredient
   * @param validIngredientId
   * @returns any
   * @throws ApiError
   */
  public static getValidIngredient(validIngredientId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredients/{validIngredientID}',
      path: {
        validIngredientID: validIngredientId,
      },
    });
  }
  /**
   * Operation for updating ValidIngredient
   * @param validIngredientId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidIngredient(
    validIngredientId: string,
    requestBody: ValidIngredientUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_ingredients/{validIngredientID}',
      path: {
        validIngredientID: validIngredientId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidIngredient
   * @param q the search query parameter
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
  public static getValidIngredientsByPreparation(
    q: string,
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
      data?: Array<ValidIngredient>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredients/by_preparation/{validPreparationID}',
      path: {
        validPreparationID: validPreparationId,
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
   * Operation for fetching ValidIngredient
   * @returns any
   * @throws ApiError
   */
  public static getRandomValidIngredient(): CancelablePromise<
    APIResponse & {
      data?: ValidIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredients/random',
    });
  }
  /**
   * Operation for fetching ValidIngredient
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
  public static searchForValidIngredients(
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
      data?: Array<ValidIngredient>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredients/search',
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
