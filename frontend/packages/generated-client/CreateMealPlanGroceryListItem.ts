// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  MealPlanGroceryListItem, 
  APIResponse, 
  MealPlanGroceryListItemCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createMealPlanGroceryListItem(
  client: Axios,
  mealPlanID: string,
  input: MealPlanGroceryListItemCreationRequestInput,
): Promise<  APIResponse <  MealPlanGroceryListItem >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < MealPlanGroceryListItem  >  >(`/api/v1/meal_plans/${mealPlanID}/grocery_list_items`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}