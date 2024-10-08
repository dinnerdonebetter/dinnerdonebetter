// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlanOption, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function getMealPlanOptions(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  mealPlanID: string,
  mealPlanEventID: string,
): Promise<QueryFilteredResult<MealPlanOption>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<MealPlanOption>>>(
      `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`,
      {
        params: filter.asRecord(),
      },
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<MealPlanOption>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
