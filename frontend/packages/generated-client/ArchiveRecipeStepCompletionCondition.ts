// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipeStepCompletionCondition, APIResponse } from '@dinnerdonebetter/models';

export async function archiveRecipeStepCompletionCondition(
  client: Axios,
  recipeID: string,
  recipeStepID: string,
  recipeStepCompletionConditionID: string,
): Promise<APIResponse<RecipeStepCompletionCondition>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<RecipeStepCompletionCondition>>(
      `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
