/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { FinalizeMealPlansRequest } from '../models/FinalizeMealPlansRequest';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class WorkersService {
  /**
   * Operation for creating FinalizeMealPlansRequest
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static runFinalizeMealPlanWorker(requestBody: FinalizeMealPlansRequest): CancelablePromise<
    APIResponse & {
      data?: FinalizeMealPlansRequest;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/workers/finalize_meal_plans',
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for creating
   * @throws ApiError
   */
  public static runMealPlanGroceryListInitializerWorker(): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/workers/meal_plan_grocery_list_init',
    });
  }
  /**
   * Operation for creating
   * @throws ApiError
   */
  public static runMealPlanTaskCreatorWorker(): CancelablePromise<void> {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/workers/meal_plan_tasks',
    });
  }
}
