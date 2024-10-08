// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlanEvent, APIResponse } from '@dinnerdonebetter/models';

export async function getMealPlanEvent(
  client: Axios,
  mealPlanID: string,
  mealPlanEventID: string,
): Promise<APIResponse<MealPlanEvent>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<MealPlanEvent>>(
      `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
