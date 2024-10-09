// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  MealPlan, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getMealPlans(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  ): Promise< QueryFilteredResult< MealPlan >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<MealPlan>  >  >(`/api/v1/meal_plans`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<MealPlan>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}