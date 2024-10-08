// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlan, APIResponse } from '@dinnerdonebetter/models';

export async function archiveMealPlan(client: Axios, mealPlanID: string): Promise<APIResponse<MealPlan>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<MealPlan>>(`/api/v1/meal_plans/${mealPlanID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
