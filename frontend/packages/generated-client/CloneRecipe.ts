// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  Recipe, 
  APIResponse, 
  
} from '@dinnerdonebetter/models';

export async function cloneRecipe(
  client: Axios,
  recipeID: string,
  
): Promise<  APIResponse <  Recipe >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < Recipe  >  >(`/api/v1/recipes/${recipeID}/clone`);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}