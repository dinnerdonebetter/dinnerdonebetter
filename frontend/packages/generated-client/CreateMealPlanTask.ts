// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlanTask, APIResponse, MealPlanTaskCreationRequestInput } from '@dinnerdonebetter/models';

export async function createMealPlanTask(
  client: Axios,
  mealPlanID: string,
  input: MealPlanTaskCreationRequestInput,
): Promise<APIResponse<MealPlanTask>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<MealPlanTask>>(`/api/v1/meal_plans/${mealPlanID}/tasks`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
