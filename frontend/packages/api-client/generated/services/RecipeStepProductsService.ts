/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { RecipeStepProduct } from '../models/RecipeStepProduct';
import type { RecipeStepProductCreationRequestInput } from '../models/RecipeStepProductCreationRequestInput';
import type { RecipeStepProductUpdateRequestInput } from '../models/RecipeStepProductUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class RecipeStepProductsService {
  /**
   * Operation for fetching RecipeStepProduct
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
  public static getRecipeStepProducts(
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
      data?: Array<RecipeStepProduct>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/products',
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
   * Operation for creating RecipeStepProduct
   * @param recipeId
   * @param recipeStepId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeStepProduct(
    recipeId: string,
    recipeStepId: string,
    requestBody: RecipeStepProductCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepProduct;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/products',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeStepProduct
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepProductId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeStepProduct(
    recipeId: string,
    recipeStepId: string,
    recipeStepProductId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepProduct;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepProductID: recipeStepProductId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStepProduct
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepProductId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepProduct(
    recipeId: string,
    recipeStepId: string,
    recipeStepProductId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepProduct;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepProductID: recipeStepProductId,
      },
    });
  }
  /**
   * Operation for updating RecipeStepProduct
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepProductId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeStepProduct(
    recipeId: string,
    recipeStepId: string,
    recipeStepProductId: string,
    requestBody: RecipeStepProductUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepProduct;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepProductID: recipeStepProductId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
