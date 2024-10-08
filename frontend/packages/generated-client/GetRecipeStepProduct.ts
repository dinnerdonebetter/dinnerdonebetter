// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipeStepProduct, APIResponse } from '@dinnerdonebetter/models';

export async function getRecipeStepProduct(
  client: Axios,
  recipeID: string,
  recipeStepID: string,
  recipeStepProductID: string,
): Promise<APIResponse<RecipeStepProduct>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<RecipeStepProduct>>(
      `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
