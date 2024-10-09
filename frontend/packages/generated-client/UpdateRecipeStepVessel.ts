// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeStepVessel, 
  APIResponse, 
  RecipeStepVesselUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateRecipeStepVessel(
  client: Axios,
  recipeID: string,recipeStepID: string,recipeStepVesselID: string,
  input: RecipeStepVesselUpdateRequestInput,
): Promise<  APIResponse <  RecipeStepVessel >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < RecipeStepVessel  >  >(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}