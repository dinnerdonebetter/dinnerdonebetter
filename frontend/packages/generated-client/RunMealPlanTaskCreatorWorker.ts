// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  CreateMealPlanTasksResponse, 
  APIResponse, 
  CreateMealPlanTasksRequest, 
} from '@dinnerdonebetter/models';

export async function runMealPlanTaskCreatorWorker(
  client: Axios,
  
  input: CreateMealPlanTasksRequest,
): Promise<  APIResponse <  CreateMealPlanTasksResponse >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < CreateMealPlanTasksResponse  >  >(`/api/v1/workers/meal_plan_tasks`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}