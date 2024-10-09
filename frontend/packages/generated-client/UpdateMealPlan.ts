// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  MealPlan, 
  APIResponse, 
  MealPlanUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateMealPlan(
  client: Axios,
  mealPlanID: string,
  input: MealPlanUpdateRequestInput,
): Promise<  APIResponse <  MealPlan >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < MealPlan  >  >(`/api/v1/meal_plans/${mealPlanID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}