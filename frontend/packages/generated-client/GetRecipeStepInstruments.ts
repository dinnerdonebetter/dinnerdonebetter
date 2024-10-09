// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeStepInstrument, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getRecipeStepInstruments(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  recipeID: string,
	recipeStepID: string,
	): Promise< QueryFilteredResult< RecipeStepInstrument >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<RecipeStepInstrument>  >  >(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<RecipeStepInstrument>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}