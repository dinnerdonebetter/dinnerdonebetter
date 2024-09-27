/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { HouseholdInstrumentOwnership } from '../models/HouseholdInstrumentOwnership';
import type { HouseholdInstrumentOwnershipCreationRequestInput } from '../models/HouseholdInstrumentOwnershipCreationRequestInput';
import type { HouseholdInstrumentOwnershipUpdateRequestInput } from '../models/HouseholdInstrumentOwnershipUpdateRequestInput';
import type { RecipeStepInstrument } from '../models/RecipeStepInstrument';
import type { RecipeStepInstrumentCreationRequestInput } from '../models/RecipeStepInstrumentCreationRequestInput';
import type { RecipeStepInstrumentUpdateRequestInput } from '../models/RecipeStepInstrumentUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class HouseholdInstrumentOwnershipsService {
  /**
   * Operation for fetching HouseholdInstrumentOwnership
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
  public static getHouseholdInstrumentOwnerships(
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
      data?: Array<HouseholdInstrumentOwnership>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/households/instruments',
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
   * Operation for creating HouseholdInstrumentOwnership
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createHouseholdInstrumentOwnership(
    requestBody: HouseholdInstrumentOwnershipCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: HouseholdInstrumentOwnership;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/households/instruments',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving HouseholdInstrumentOwnership
   * @param householdInstrumentOwnershipId
   * @returns any
   * @throws ApiError
   */
  public static archiveHouseholdInstrumentOwnership(householdInstrumentOwnershipId: string): CancelablePromise<
    APIResponse & {
      data?: HouseholdInstrumentOwnership;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/households/instruments/{householdInstrumentOwnershipID}',
      path: {
        householdInstrumentOwnershipID: householdInstrumentOwnershipId,
      },
    });
  }
  /**
   * Operation for fetching HouseholdInstrumentOwnership
   * @param householdInstrumentOwnershipId
   * @returns any
   * @throws ApiError
   */
  public static getHouseholdInstrumentOwnership(householdInstrumentOwnershipId: string): CancelablePromise<
    APIResponse & {
      data?: HouseholdInstrumentOwnership;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/households/instruments/{householdInstrumentOwnershipID}',
      path: {
        householdInstrumentOwnershipID: householdInstrumentOwnershipId,
      },
    });
  }
  /**
   * Operation for updating HouseholdInstrumentOwnership
   * @param householdInstrumentOwnershipId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateHouseholdInstrumentOwnership(
    householdInstrumentOwnershipId: string,
    requestBody: HouseholdInstrumentOwnershipUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: HouseholdInstrumentOwnership;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/households/instruments/{householdInstrumentOwnershipID}',
      path: {
        householdInstrumentOwnershipID: householdInstrumentOwnershipId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching RecipeStepInstrument
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
  public static getRecipeStepInstruments(
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
      data?: Array<RecipeStepInstrument>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments',
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
   * Operation for creating RecipeStepInstrument
   * @param recipeId
   * @param recipeStepId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createRecipeStepInstrument(
    recipeId: string,
    recipeStepId: string,
    requestBody: RecipeStepInstrumentCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving RecipeStepInstrument
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepInstrumentId
   * @returns any
   * @throws ApiError
   */
  public static archiveRecipeStepInstrument(
    recipeId: string,
    recipeStepId: string,
    recipeStepInstrumentId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepInstrumentID: recipeStepInstrumentId,
      },
    });
  }
  /**
   * Operation for fetching RecipeStepInstrument
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepInstrumentId
   * @returns any
   * @throws ApiError
   */
  public static getRecipeStepInstrument(
    recipeId: string,
    recipeStepId: string,
    recipeStepInstrumentId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepInstrumentID: recipeStepInstrumentId,
      },
    });
  }
  /**
   * Operation for updating RecipeStepInstrument
   * @param recipeId
   * @param recipeStepId
   * @param recipeStepInstrumentId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateRecipeStepInstrument(
    recipeId: string,
    recipeStepId: string,
    recipeStepInstrumentId: string,
    requestBody: RecipeStepInstrumentUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: RecipeStepInstrument;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}',
      path: {
        recipeID: recipeId,
        recipeStepID: recipeStepId,
        recipeStepInstrumentID: recipeStepInstrumentId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
