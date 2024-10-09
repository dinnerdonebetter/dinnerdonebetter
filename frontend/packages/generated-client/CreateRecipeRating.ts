// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeRating, 
  APIResponse, 
  RecipeRatingCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createRecipeRating(
  client: Axios,
  recipeID: string,
  input: RecipeRatingCreationRequestInput,
): Promise<  APIResponse <  RecipeRating >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < RecipeRating  >  >(`/api/v1/recipes/${recipeID}/ratings`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}