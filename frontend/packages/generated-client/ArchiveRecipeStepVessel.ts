// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipeStepVessel, APIResponse } from '@dinnerdonebetter/models';

export async function archiveRecipeStepVessel(
  client: Axios,
  recipeID: string,
  recipeStepID: string,
  recipeStepVesselID: string,
): Promise<APIResponse<RecipeStepVessel>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<RecipeStepVessel>>(
      `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
