// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeStepVessel, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getRecipeStepVessels(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  recipeID: string,
	recipeStepID: string,
	): Promise< QueryFilteredResult< RecipeStepVessel >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<RecipeStepVessel>  >  >(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<RecipeStepVessel>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}