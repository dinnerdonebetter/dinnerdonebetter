// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  RecipeStepCompletionCondition,
  APIResponse,
  RecipeStepCompletionConditionForExistingRecipeCreationRequestInput,
} from '@dinnerdonebetter/models';

export async function createRecipeStepCompletionCondition(
  client: Axios,
  recipeID: string,
  recipeStepID: string,
  input: RecipeStepCompletionConditionForExistingRecipeCreationRequestInput,
): Promise<APIResponse<RecipeStepCompletionCondition>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<RecipeStepCompletionCondition>>(
      `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
