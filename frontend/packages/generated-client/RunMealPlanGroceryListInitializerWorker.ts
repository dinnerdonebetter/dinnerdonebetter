// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  InitializeMealPlanGroceryListResponse,
  APIResponse,
  InitializeMealPlanGroceryListRequest,
} from '@dinnerdonebetter/models';

export async function runMealPlanGroceryListInitializerWorker(
  client: Axios,

  input: InitializeMealPlanGroceryListRequest,
): Promise<APIResponse<InitializeMealPlanGroceryListResponse>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<InitializeMealPlanGroceryListResponse>>(
      `/api/v1/workers/meal_plan_grocery_list_init`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
