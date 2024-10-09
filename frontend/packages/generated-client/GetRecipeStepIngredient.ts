// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeStepIngredient, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getRecipeStepIngredient(
  client: Axios,
  recipeID: string,
	recipeStepID: string,
	recipeStepIngredientID: string,
	): Promise<  APIResponse <  RecipeStepIngredient >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < RecipeStepIngredient  >  >(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}