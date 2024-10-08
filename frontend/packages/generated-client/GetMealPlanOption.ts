// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlanOption, APIResponse } from '@dinnerdonebetter/models';

export async function getMealPlanOption(
  client: Axios,
  mealPlanID: string,
  mealPlanEventID: string,
  mealPlanOptionID: string,
): Promise<APIResponse<MealPlanOption>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<MealPlanOption>>(
      `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
