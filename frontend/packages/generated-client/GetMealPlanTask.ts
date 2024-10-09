// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  MealPlanTask, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getMealPlanTask(
  client: Axios,
  mealPlanID: string,
	mealPlanTaskID: string,
	): Promise<  APIResponse <  MealPlanTask >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < MealPlanTask  >  >(`/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}