/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { RecipeRating } from '../models/RecipeRating';
import type { RecipeRatingCreationRequestInput } from '../models/RecipeRatingCreationRequestInput';
import type { RecipeRatingUpdateRequestInput } from '../models/RecipeRatingUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class RecipeRatingsService {
  /**
   * Operation for fetching RecipeRating
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param recipeId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeRatings(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    recipeId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<RecipeRating>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/ratings',
      path: {
        recipeID: recipeId,
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
   * Operation for creating RecipeRating
   * @param recipeId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeRating(
    recipeId: string,
    requestBody: RecipeRatingCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeRating;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/ratings',
      path: {
        recipeID: recipeId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeRating
   * @param recipeId
   * @param recipeRatingId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeRating(
    recipeId: string,
    recipeRatingId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeRating;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/ratings/{recipeRatingID}',
      path: {
        recipeID: recipeId,
        recipeRatingID: recipeRatingId,
      },
    });
  }
  /**
   * Operation for fetching RecipeRating
   * @param recipeId
   * @param recipeRatingId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeRating(
    recipeId: string,
    recipeRatingId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeRating;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/ratings/{recipeRatingID}',
      path: {
        recipeID: recipeId,
        recipeRatingID: recipeRatingId,
      },
    });
  }
  /**
   * Operation for updating RecipeRating
   * @param recipeId
   * @param recipeRatingId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeRating(
    recipeId: string,
    recipeRatingId: string,
    requestBody: RecipeRatingUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeRating;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/ratings/{recipeRatingID}',
      path: {
        recipeID: recipeId,
        recipeRatingID: recipeRatingId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
