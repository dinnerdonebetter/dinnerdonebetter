/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { MealPlanTask } from '../models/MealPlanTask';
import type { MealPlanTaskCreationRequestInput } from '../models/MealPlanTaskCreationRequestInput';
import type { MealPlanTaskStatusChangeRequestInput } from '../models/MealPlanTaskStatusChangeRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class MealPlanTasksService {
  /**
   * Operation for fetching MealPlanTask
   * @param limit How many results should appear in output, max is 250.
   * @param page What page of results should appear in output.
   * @param createdBefore The latest CreatedAt date that should appear in output.
   * @param createdAfter The earliest CreatedAt date that should appear in output.
   * @param updatedBefore The latest UpdatedAt date that should appear in output.
   * @param updatedAfter The earliest UpdatedAt date that should appear in output.
   * @param includeArchived Whether or not to include archived results in output, limited to service admins.
   * @param sortBy The direction in which results should be sorted.
   * @param mealPlanId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlanTasks(
    limit: number,
    page: number,
    createdBefore: string,
    createdAfter: string,
    updatedBefore: string,
    updatedAfter: string,
    includeArchived: 'true' | 'false',
    sortBy: 'asc' | 'desc',
    mealPlanId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: Array<MealPlanTask>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/tasks',
      path: {
        mealPlanID: mealPlanId,
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
   * Operation for creating MealPlanTask
   * @param mealPlanId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createMealPlanTask(
    mealPlanId: string,
    requestBody: MealPlanTaskCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanTask;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/meal_plans/{mealPlanID}/tasks',
      path: {
        mealPlanID: mealPlanId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for fetching MealPlanTask
   * @param mealPlanId
   * @param mealPlanTaskId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlanTask(
    mealPlanId: string,
    mealPlanTaskId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanTask;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanTaskID: mealPlanTaskId,
      },
    });
  }
  /**
   * Operation for updating MealPlanTask
   * @param mealPlanId
   * @param mealPlanTaskId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateMealPlanTaskStatus(
    mealPlanId: string,
    mealPlanTaskId: string,
    requestBody: MealPlanTaskStatusChangeRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanTask;
    }
  > {
    return __request(OpenAPI, {
      method: 'PATCH',
      url: '/api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanTaskID: mealPlanTaskId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
