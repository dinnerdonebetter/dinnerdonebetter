// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlan, APIResponse, MealPlanCreationRequestInput } from '@dinnerdonebetter/models';

export async function createMealPlan(
  client: Axios,

  input: MealPlanCreationRequestInput,
): Promise<APIResponse<MealPlan>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<MealPlan>>(`/api/v1/meal_plans`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
