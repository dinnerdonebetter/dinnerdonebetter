/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidIngredientStateIngredient } from '../models/ValidIngredientStateIngredient';
import type { ValidIngredientStateIngredientCreationRequestInput } from '../models/ValidIngredientStateIngredientCreationRequestInput';
import type { ValidIngredientStateIngredientUpdateRequestInput } from '../models/ValidIngredientStateIngredientUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidIngredientStateIngredientsService {
  /**
   * Operation for fetching ValidIngredientStateIngredient
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
  public static getValidIngredientStateIngredients(
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
      data?: Array<ValidIngredientStateIngredient>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_state_ingredients',
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
   * Operation for creating ValidIngredientStateIngredient
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidIngredientStateIngredient(
    requestBody: ValidIngredientStateIngredientCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientStateIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_ingredient_state_ingredients',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidIngredientStateIngredient
   * @param validIngredientStateIngredientId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidIngredientStateIngredient(validIngredientStateIngredientId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientStateIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}',
      path: {
        validIngredientStateIngredientID: validIngredientStateIngredientId,
      },
    });
  }
  /**
   * Operation for fetching ValidIngredientStateIngredient
   * @param validIngredientStateIngredientId
   * @returns any
   * @throws ApiError
   */
  public static getValidIngredientStateIngredient(validIngredientStateIngredientId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientStateIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}',
      path: {
        validIngredientStateIngredientID: validIngredientStateIngredientId,
      },
    });
  }
  /**
   * Operation for updating ValidIngredientStateIngredient
   * @param validIngredientStateIngredientId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidIngredientStateIngredient(
    validIngredientStateIngredientId: string,
    requestBody: ValidIngredientStateIngredientUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientStateIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}',
      path: {
        validIngredientStateIngredientID: validIngredientStateIngredientId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidIngredientStateIngredient
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
  public static getValidIngredientStateIngredientsByIngredient(
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
      data?: Array<ValidIngredientStateIngredient>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_state_ingredients/by_ingredient/{validIngredientID}',
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
   * Operation for fetching ValidIngredientStateIngredient
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param validIngredientStateId
   * @returns any
   * @throws ApiError
   */
  public static getValidIngredientStateIngredientsByIngredientState(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    validIngredientStateId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<ValidIngredientStateIngredient>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/{validIngredientStateID}',
      path: {
        validIngredientStateID: validIngredientStateId,
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
