// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeStep, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function archiveRecipeStep(
  client: Axios,
  recipeID: string,
	recipeStepID: string,
	): Promise<  APIResponse <  RecipeStep >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete< APIResponse < RecipeStep  >  >(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}