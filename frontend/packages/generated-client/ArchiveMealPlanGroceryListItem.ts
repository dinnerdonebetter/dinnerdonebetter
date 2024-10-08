// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlanGroceryListItem, APIResponse } from '@dinnerdonebetter/models';

export async function archiveMealPlanGroceryListItem(
  client: Axios,
  mealPlanID: string,
  mealPlanGroceryListItemID: string,
): Promise<APIResponse<MealPlanGroceryListItem>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<MealPlanGroceryListItem>>(
      `/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
