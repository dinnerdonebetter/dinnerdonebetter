// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlanGroceryListItem, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function getMealPlanGroceryListItemsForMealPlan(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  mealPlanID: string,
): Promise<QueryFilteredResult<MealPlanGroceryListItem>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<MealPlanGroceryListItem>>>(
      `/api/v1/meal_plans/${mealPlanID}/grocery_list_items`,
      {
        params: filter.asRecord(),
      },
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<MealPlanGroceryListItem>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
