// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  MealPlanEvent, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getMealPlanEvents(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  mealPlanID: string,
	): Promise< QueryFilteredResult< MealPlanEvent >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<MealPlanEvent>  >  >(`/api/v1/meal_plans/${mealPlanID}/events`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<MealPlanEvent>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}