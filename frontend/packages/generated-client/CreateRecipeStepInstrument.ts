// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipeStepInstrument, APIResponse, RecipeStepInstrumentCreationRequestInput } from '@dinnerdonebetter/models';

export async function createRecipeStepInstrument(
  client: Axios,
  recipeID: string,
  recipeStepID: string,
  input: RecipeStepInstrumentCreationRequestInput,
): Promise<APIResponse<RecipeStepInstrument>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<RecipeStepInstrument>>(
      `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
