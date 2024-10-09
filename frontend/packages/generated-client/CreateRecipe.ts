// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  Recipe, 
  APIResponse, 
  RecipeCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createRecipe(
  client: Axios,
  
  input: RecipeCreationRequestInput,
): Promise<  APIResponse <  Recipe >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < Recipe  >  >(`/api/v1/recipes`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}