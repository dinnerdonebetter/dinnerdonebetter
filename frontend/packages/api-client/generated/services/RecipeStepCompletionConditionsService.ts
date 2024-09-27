/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { RecipeStepCompletionCondition } from '../models/RecipeStepCompletionCondition';
import type { RecipeStepCompletionConditionForExistingRecipeCreationRequestInput } from '../models/RecipeStepCompletionConditionForExistingRecipeCreationRequestInput';
import type { RecipeStepCompletionConditionUpdateRequestInput } from '../models/RecipeStepCompletionConditionUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class RecipeStepCompletionConditionsService {
  /**
   * Operation for fetching RecipeStepCompletionCondition
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param recipeId
   * @param recipeStepId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepCompletionConditions(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    recipeId: string,
    recipeStepId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<RecipeStepCompletionCondition>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
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
   * Operation for creating RecipeStepCompletionCondition
   * @param recipeId
   * @param recipeStepId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeStepCompletionCondition(
    recipeId: string,
    recipeStepId: string,
    requestBody: RecipeStepCompletionConditionForExistingRecipeCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepCompletionCondition;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeStepCompletionCondition
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepCompletionConditionId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeStepCompletionCondition(
    recipeId: string,
    recipeStepId: string,
    recipeStepCompletionConditionId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepCompletionCondition;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepCompletionConditionID: recipeStepCompletionConditionId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStepCompletionCondition
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepCompletionConditionId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepCompletionCondition(
    recipeId: string,
    recipeStepId: string,
    recipeStepCompletionConditionId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepCompletionCondition;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepCompletionConditionID: recipeStepCompletionConditionId,
      },
    });
  }
  /**
   * Operation for updating RecipeStepCompletionCondition
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepCompletionConditionId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeStepCompletionCondition(
    recipeId: string,
    recipeStepId: string,
    recipeStepCompletionConditionId: string,
    requestBody: RecipeStepCompletionConditionUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepCompletionCondition;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepCompletionConditionID: recipeStepCompletionConditionId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
