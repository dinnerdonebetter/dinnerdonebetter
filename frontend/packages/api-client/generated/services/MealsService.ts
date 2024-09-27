/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { Meal } from '../models/Meal';
import type { MealCreationRequestInput } from '../models/MealCreationRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class MealsService {
  /**
   * Operation for fetching Meal
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
  public static getMeals(
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
      data?: Array<Meal>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meals',
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
   * Operation for creating Meal
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createMeal(requestBody: MealCreationRequestInput): CancelablePromise<
    APIResponse & {
      data?: Meal;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/meals',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving Meal
   * @param mealId
   * @returns any
   * @throws ApiError
   */
  public static archiveMeal(mealId: string): CancelablePromise<
    APIResponse & {
      data?: Meal;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/meals/{mealID}',
      path: {
        mealID: mealId,
      },
    });
  }
  /**
   * Operation for fetching Meal
   * @param mealId
   * @returns any
   * @throws ApiError
   */
  public static getMeal(mealId: string): CancelablePromise<
    APIResponse & {
      data?: Meal;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meals/{mealID}',
      path: {
        mealID: mealId,
      },
    });
  }
  /**
   * Operation for fetching Meal
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
  public static searchForMeals(
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
      data?: Array<Meal>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meals/search',
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
