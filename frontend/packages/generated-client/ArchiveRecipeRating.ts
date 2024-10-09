// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeRating, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function archiveRecipeRating(
  client: Axios,
  recipeID: string,
	recipeRatingID: string,
	): Promise<  APIResponse <  RecipeRating >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete< APIResponse < RecipeRating  >  >(`/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}