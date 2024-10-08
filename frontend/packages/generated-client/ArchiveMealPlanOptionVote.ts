// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlanOptionVote, APIResponse } from '@dinnerdonebetter/models';

export async function archiveMealPlanOptionVote(
  client: Axios,
  mealPlanID: string,
  mealPlanEventID: string,
  mealPlanOptionID: string,
  mealPlanOptionVoteID: string,
): Promise<APIResponse<MealPlanOptionVote>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<MealPlanOptionVote>>(
      `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
