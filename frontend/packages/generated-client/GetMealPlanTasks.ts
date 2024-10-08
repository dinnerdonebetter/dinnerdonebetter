// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { MealPlanTask, QueryFilter, QueryFilteredResult, APIResponse } from '@dinnerdonebetter/models';

export async function getMealPlanTasks(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  mealPlanID: string,
): Promise<QueryFilteredResult<MealPlanTask>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Array<MealPlanTask>>>(`/api/v1/meal_plans/${mealPlanID}/tasks`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<MealPlanTask>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
