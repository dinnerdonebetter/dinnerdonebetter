// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlanOption, APIResponse, MealPlanOptionUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateMealPlanOption(
  client: Axios,
  mealPlanID: string,
  mealPlanEventID: string,
  mealPlanOptionID: string,
  input: MealPlanOptionUpdateRequestInput,
): Promise<APIResponse<MealPlanOption>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<MealPlanOption>>(
      `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
