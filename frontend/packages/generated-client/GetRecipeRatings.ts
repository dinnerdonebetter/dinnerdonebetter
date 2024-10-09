// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeRating, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getRecipeRatings(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  recipeID: string,
	): Promise< QueryFilteredResult< RecipeRating >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<RecipeRating>  >  >(`/api/v1/recipes/${recipeID}/ratings`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<RecipeRating>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}