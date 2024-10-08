// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipeStepProduct, APIResponse, RecipeStepProductUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateRecipeStepProduct(
  client: Axios,
  recipeID: string,
  recipeStepID: string,
  recipeStepProductID: string,
  input: RecipeStepProductUpdateRequestInput,
): Promise<APIResponse<RecipeStepProduct>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<RecipeStepProduct>>(
      `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
