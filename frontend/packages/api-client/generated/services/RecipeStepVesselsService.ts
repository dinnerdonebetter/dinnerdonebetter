/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { RecipeStepVessel } from '../models/RecipeStepVessel';
import type { RecipeStepVesselCreationRequestInput } from '../models/RecipeStepVesselCreationRequestInput';
import type { RecipeStepVesselUpdateRequestInput } from '../models/RecipeStepVesselUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class RecipeStepVesselsService {
  /**
   * Operation for fetching RecipeStepVessel
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
  public static getRecipeStepVessels(
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
      data?: Array<RecipeStepVessel>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels',
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
   * Operation for creating RecipeStepVessel
   * @param recipeId
   * @param recipeStepId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeStepVessel(
    recipeId: string,
    recipeStepId: string,
    requestBody: RecipeStepVesselCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeStepVessel
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepVesselId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeStepVessel(
    recipeId: string,
    recipeStepId: string,
    recipeStepVesselId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepVesselID: recipeStepVesselId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStepVessel
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepVesselId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepVessel(
    recipeId: string,
    recipeStepId: string,
    recipeStepVesselId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepVesselID: recipeStepVesselId,
      },
    });
  }
  /**
   * Operation for updating RecipeStepVessel
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepVesselId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeStepVessel(
    recipeId: string,
    recipeStepId: string,
    recipeStepVesselId: string,
    requestBody: RecipeStepVesselUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepVessel;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepVesselID: recipeStepVesselId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
