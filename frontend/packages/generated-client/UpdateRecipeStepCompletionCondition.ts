// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeStepCompletionCondition, 
  APIResponse, 
  RecipeStepCompletionConditionUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateRecipeStepCompletionCondition(
  client: Axios,
  recipeID: string,recipeStepID: string,recipeStepCompletionConditionID: string,
  input: RecipeStepCompletionConditionUpdateRequestInput,
): Promise<  APIResponse <  RecipeStepCompletionCondition >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < RecipeStepCompletionCondition  >  >(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}