// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  MealPlanOptionVote, 
  APIResponse, 
  MealPlanOptionVoteCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createMealPlanOptionVote(
  client: Axios,
  mealPlanID: string,mealPlanEventID: string,
  input: MealPlanOptionVoteCreationRequestInput,
): Promise<  APIResponse <  MealPlanOptionVote >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < MealPlanOptionVote  >  >(`/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/vote`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}