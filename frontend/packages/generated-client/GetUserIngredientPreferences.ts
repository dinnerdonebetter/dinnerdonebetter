// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  UserIngredientPreference, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getUserIngredientPreferences(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  ): Promise< QueryFilteredResult< UserIngredientPreference >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<UserIngredientPreference>  >  >(`/api/v1/user_ingredient_preferences`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<UserIngredientPreference>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}