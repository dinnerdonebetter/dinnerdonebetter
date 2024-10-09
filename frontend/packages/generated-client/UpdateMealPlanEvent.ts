// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  MealPlanEvent, 
  APIResponse, 
  MealPlanEventUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateMealPlanEvent(
  client: Axios,
  mealPlanID: string,mealPlanEventID: string,
  input: MealPlanEventUpdateRequestInput,
): Promise<  APIResponse <  MealPlanEvent >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < MealPlanEvent  >  >(`/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}