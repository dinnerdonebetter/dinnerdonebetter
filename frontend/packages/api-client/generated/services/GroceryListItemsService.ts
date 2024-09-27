/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { APIResponse } from '../models/APIResponse';
import type { MealPlanGroceryListItem } from '../models/MealPlanGroceryListItem';
import type { MealPlanGroceryListItemCreationRequestInput } from '../models/MealPlanGroceryListItemCreationRequestInput';
import type { MealPlanGroceryListItemUpdateRequestInput } from '../models/MealPlanGroceryListItemUpdateRequestInput';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class GroceryListItemsService {
  /**
   * Operation for fetching MealPlanGroceryListItem
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
  public static getMealPlanGroceryListItemsForMealPlan(
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
      data?: Array<MealPlanGroceryListItem>;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/grocery_list_items',
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
   * Operation for creating MealPlanGroceryListItem
   * @param mealPlanId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static createMealPlanGroceryListItem(
    mealPlanId: string,
    requestBody: MealPlanGroceryListItemCreationRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanGroceryListItem;
    }
  > {
    return __request(OpenAPI, {
      method: 'POST',
      url: '/api/v1/meal_plans/{mealPlanID}/grocery_list_items',
      path: {
        mealPlanID: mealPlanId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
  /**
   * Operation for archiving MealPlanGroceryListItem
   * @param mealPlanId
   * @param mealPlanGroceryListItemId
   * @returns any
   * @throws ApiError
   */
  public static archiveMealPlanGroceryListItem(
    mealPlanId: string,
    mealPlanGroceryListItemId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanGroceryListItem;
    }
  > {
    return __request(OpenAPI, {
      method: 'DELETE',
      url: '/api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanGroceryListItemID: mealPlanGroceryListItemId,
      },
    });
  }
  /**
   * Operation for fetching MealPlanGroceryListItem
   * @param mealPlanId
   * @param mealPlanGroceryListItemId
   * @returns any
   * @throws ApiError
   */
  public static getMealPlanGroceryListItem(
    mealPlanId: string,
    mealPlanGroceryListItemId: string,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanGroceryListItem;
    }
  > {
    return __request(OpenAPI, {
      method: 'GET',
      url: '/api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanGroceryListItemID: mealPlanGroceryListItemId,
      },
    });
  }
  /**
   * Operation for updating MealPlanGroceryListItem
   * @param mealPlanId
   * @param mealPlanGroceryListItemId
   * @param requestBody
   * @returns any
   * @throws ApiError
   */
  public static updateMealPlanGroceryListItem(
    mealPlanId: string,
    mealPlanGroceryListItemId: string,
    requestBody: MealPlanGroceryListItemUpdateRequestInput,
  ): CancelablePromise<
    APIResponse & {
      data?: MealPlanGroceryListItem;
    }
  > {
    return __request(OpenAPI, {
      method: 'PUT',
      url: '/api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}',
      path: {
        mealPlanID: mealPlanId,
        mealPlanGroceryListItemID: mealPlanGroceryListItemId,
      },
      body: requestBody,
      mediaType: 'application/json',
    });
  }
}
