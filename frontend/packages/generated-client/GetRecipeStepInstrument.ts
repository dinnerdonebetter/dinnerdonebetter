// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeStepInstrument, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getRecipeStepInstrument(
  client: Axios,
  recipeID: string,
	recipeStepID: string,
	recipeStepInstrumentID: string,
	): Promise<  APIResponse <  RecipeStepInstrument >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < RecipeStepInstrument  >  >(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}