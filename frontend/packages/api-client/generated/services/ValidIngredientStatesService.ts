/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidIngredientState } from '../models/ValidIngredientState';
import type { ValidIngredientStateCreationRequestInput } from '../models/ValidIngredientStateCreationRequestInput';
import type { ValidIngredientStateUpdateRequestInput } from '../models/ValidIngredientStateUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidIngredientStatesService {
  /**
   * Operation for fetching ValidIngredientState
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
  public static getValidIngredientStates(
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
      data?: Array<ValidIngredientState>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_states',
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
   * Operation for creating ValidIngredientState
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidIngredientState(requestBody: ValidIngredientStateCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientState;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_ingredient_states',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidIngredientState
   * @param validIngredientStateId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidIngredientState(validIngredientStateId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientState;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_ingredient_states/{validIngredientStateID}',
      path: {
        validIngredientStateID: validIngredientStateId,
      },
    });
  }
  /**
   * Operation for fetching ValidIngredientState
   * @param validIngredientStateId
   * @returns any
   * @throws ApiError
   */
  public static getValidIngredientState(validIngredientStateId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientState;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_states/{validIngredientStateID}',
      path: {
        validIngredientStateID: validIngredientStateId,
      },
    });
  }
  /**
   * Operation for updating ValidIngredientState
   * @param validIngredientStateId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidIngredientState(
    validIngredientStateId: string,
    requestBody: ValidIngredientStateUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientState;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_ingredient_states/{validIngredientStateID}',
      path: {
        validIngredientStateID: validIngredientStateId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidIngredientState
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
  public static searchForValidIngredientStates(
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
      data?: Array<ValidIngredientState>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_states/search',
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
