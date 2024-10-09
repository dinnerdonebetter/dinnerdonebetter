// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  FinalizeMealPlansResponse, 
  APIResponse, 
  
} from '@dinnerdonebetter/models';

export async function finalizeMealPlan(
  client: Axios,
  mealPlanID: string,
  
): Promise<  APIResponse <  FinalizeMealPlansResponse >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < FinalizeMealPlansResponse  >  >(`/api/v1/meal_plans/${mealPlanID}/finalize`);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}