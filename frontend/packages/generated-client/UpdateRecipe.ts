// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  Recipe, 
  APIResponse, 
  RecipeUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateRecipe(
  client: Axios,
  recipeID: string,
  input: RecipeUpdateRequestInput,
): Promise<  APIResponse <  Recipe >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < Recipe  >  >(`/api/v1/recipes/${recipeID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}