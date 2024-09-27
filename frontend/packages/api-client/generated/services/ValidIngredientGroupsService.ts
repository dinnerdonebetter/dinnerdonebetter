/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { ValidIngredientGroup } from '../models/ValidIngredientGroup';
import type { ValidIngredientGroupCreationRequestInput } from '../models/ValidIngredientGroupCreationRequestInput';
import type { ValidIngredientGroupUpdateRequestInput } from '../models/ValidIngredientGroupUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class ValidIngredientGroupsService {
  /**
   * Operation for fetching ValidIngredientGroup
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
  public static getValidIngredientGroups(
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
      data?: Array<ValidIngredientGroup>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_groups',
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
   * Operation for creating ValidIngredientGroup
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createValidIngredientGroup(requestBody: ValidIngredientGroupCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientGroup;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/valid_ingredient_groups',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving ValidIngredientGroup
   * @param validIngredientGroupId
   * @returns any
   * @throws ApiError
   */
  public static archiveValidIngredientGroup(validIngredientGroupId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientGroup;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/valid_ingredient_groups/{validIngredientGroupID}',
      path: {
        validIngredientGroupID: validIngredientGroupId,
      },
    });
  }
  /**
   * Operation for fetching ValidIngredientGroup
   * @param validIngredientGroupId
   * @returns any
   * @throws ApiError
   */
  public static getValidIngredientGroup(validIngredientGroupId: string): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientGroup;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_groups/{validIngredientGroupID}',
      path: {
        validIngredientGroupID: validIngredientGroupId,
      },
    });
  }
  /**
   * Operation for updating ValidIngredientGroup
   * @param validIngredientGroupId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateValidIngredientGroup(
    validIngredientGroupId: string,
    requestBody: ValidIngredientGroupUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: ValidIngredientGroup;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/valid_ingredient_groups/{validIngredientGroupID}',
      path: {
        validIngredientGroupID: validIngredientGroupId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching ValidIngredientGroup
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
  public static searchForValidIngredientGroups(
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
      data?: Array<ValidIngredientGroup>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/valid_ingredient_groups/search',
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
