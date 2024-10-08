// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipeStep, APIResponse, RecipeStepCreationRequestInput } from '@dinnerdonebetter/models';

export async function createRecipeStep(
  client: Axios,
  recipeID: string,
  input: RecipeStepCreationRequestInput,
): Promise<APIResponse<RecipeStep>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<RecipeStep>>(`/api/v1/recipes/${recipeID}/steps`, input);
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
