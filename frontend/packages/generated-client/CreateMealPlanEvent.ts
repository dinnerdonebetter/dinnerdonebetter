// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlanEvent, APIResponse, MealPlanEventCreationRequestInput } from '@dinnerdonebetter/models';

export async function createMealPlanEvent(
  client: Axios,
  mealPlanID: string,
  input: MealPlanEventCreationRequestInput,
): Promise<APIResponse<MealPlanEvent>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<MealPlanEvent>>(`/api/v1/meal_plans/${mealPlanID}/events`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
