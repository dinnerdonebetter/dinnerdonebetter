// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeStep, 
  APIResponse, 
  RecipeStepUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateRecipeStep(
  client: Axios,
  recipeID: string,recipeStepID: string,
  input: RecipeStepUpdateRequestInput,
): Promise<  APIResponse <  RecipeStep >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < RecipeStep  >  >(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}