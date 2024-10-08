// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  MealPlanGroceryListItem,
  APIResponse,
  MealPlanGroceryListItemUpdateRequestInput,
} from '@dinnerdonebetter/models';

export async function updateMealPlanGroceryListItem(
  client: Axios,
  mealPlanID: string,
  mealPlanGroceryListItemID: string,
  input: MealPlanGroceryListItemUpdateRequestInput,
): Promise<APIResponse<MealPlanGroceryListItem>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<MealPlanGroceryListItem>>(
      `/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
