// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeStepIngredient, 
  APIResponse, 
  RecipeStepIngredientCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createRecipeStepIngredient(
  client: Axios,
  recipeID: string,recipeStepID: string,
  input: RecipeStepIngredientCreationRequestInput,
): Promise<  APIResponse <  RecipeStepIngredient >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < RecipeStepIngredient  >  >(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}