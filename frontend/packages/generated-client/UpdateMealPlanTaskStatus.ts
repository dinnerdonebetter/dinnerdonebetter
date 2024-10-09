// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  MealPlanTask, 
  APIResponse, 
  MealPlanTaskStatusChangeRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateMealPlanTaskStatus(
  client: Axios,
  mealPlanID: string,mealPlanTaskID: string,
  input: MealPlanTaskStatusChangeRequestInput,
): Promise<  APIResponse <  MealPlanTask >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.patch<APIResponse < MealPlanTask  >  >(`/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}