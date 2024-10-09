// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  MealPlanOptionVote, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getMealPlanOptionVotes(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  mealPlanID: string,
	mealPlanEventID: string,
	mealPlanOptionID: string,
	): Promise< QueryFilteredResult< MealPlanOptionVote >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<MealPlanOptionVote>  >  >(`/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<MealPlanOptionVote>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}