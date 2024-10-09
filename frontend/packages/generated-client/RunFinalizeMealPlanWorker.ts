// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  FinalizeMealPlansResponse, 
  APIResponse, 
  FinalizeMealPlansRequest, 
} from '@dinnerdonebetter/models';

export async function runFinalizeMealPlanWorker(
  client: Axios,
  
  input: FinalizeMealPlansRequest,
): Promise<  APIResponse <  FinalizeMealPlansResponse >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < FinalizeMealPlansResponse  >  >(`/api/v1/workers/finalize_meal_plans`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}