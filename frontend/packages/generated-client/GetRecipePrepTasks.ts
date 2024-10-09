// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipePrepTask, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getRecipePrepTasks(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  recipeID: string,
	): Promise< QueryFilteredResult< RecipePrepTask >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<RecipePrepTask>  >  >(`/api/v1/recipes/${recipeID}/prep_tasks`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<RecipePrepTask>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}