// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipePrepTaskStep, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getRecipeMealPlanTasks(
  client: Axios,
  recipeID: string,
	): Promise<  APIResponse <  RecipePrepTaskStep >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < RecipePrepTaskStep  >  >(`/api/v1/recipes/${recipeID}/prep_steps`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}