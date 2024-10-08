// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipeStepVessel, APIResponse, RecipeStepVesselCreationRequestInput } from '@dinnerdonebetter/models';

export async function createRecipeStepVessel(
  client: Axios,
  recipeID: string,
  recipeStepID: string,
  input: RecipeStepVesselCreationRequestInput,
): Promise<APIResponse<RecipeStepVessel>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<RecipeStepVessel>>(
      `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
