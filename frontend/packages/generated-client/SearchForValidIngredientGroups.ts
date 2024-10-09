// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientGroup, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function searchForValidIngredientGroups(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  q: string,
	): Promise< QueryFilteredResult< ValidIngredientGroup >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<ValidIngredientGroup>  >  >(`/api/v1/valid_ingredient_groups/search`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidIngredientGroup>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}