// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipeStepInstrument, APIResponse, RecipeStepInstrumentUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateRecipeStepInstrument(
  client: Axios,
  recipeID: string,
  recipeStepID: string,
  recipeStepInstrumentID: string,
  input: RecipeStepInstrumentUpdateRequestInput,
): Promise<APIResponse<RecipeStepInstrument>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<RecipeStepInstrument>>(
      `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
