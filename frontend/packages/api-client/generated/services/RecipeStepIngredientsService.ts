/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { RecipeStepIngredient } from '../models/RecipeStepIngredient';
import type { RecipeStepIngredientCreationRequestInput } from '../models/RecipeStepIngredientCreationRequestInput';
import type { RecipeStepIngredientUpdateRequestInput } from '../models/RecipeStepIngredientUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class RecipeStepIngredientsService {
  /**
   * Operation for fetching RecipeStepIngredient
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
  public static getRecipeStepIngredients(
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
      data?: Array<RecipeStepIngredient>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients',
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
   * Operation for creating RecipeStepIngredient
   * @param recipeId
   * @param recipeStepId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeStepIngredient(
    recipeId: string,
    recipeStepId: string,
    requestBody: RecipeStepIngredientCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeStepIngredient
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepIngredientId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeStepIngredient(
    recipeId: string,
    recipeStepId: string,
    recipeStepIngredientId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepIngredientID: recipeStepIngredientId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStepIngredient
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepIngredientId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepIngredient(
    recipeId: string,
    recipeStepId: string,
    recipeStepIngredientId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepIngredientID: recipeStepIngredientId,
      },
    });
  }
  /**
   * Operation for updating RecipeStepIngredient
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepIngredientId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeStepIngredient(
    recipeId: string,
    recipeStepId: string,
    recipeStepIngredientId: string,
    requestBody: RecipeStepIngredientUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepIngredient;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepIngredientID: recipeStepIngredientId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
