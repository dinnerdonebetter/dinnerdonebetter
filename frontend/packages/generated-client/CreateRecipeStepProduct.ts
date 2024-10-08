// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipeStepProduct, APIResponse, RecipeStepProductCreationRequestInput } from '@dinnerdonebetter/models';

export async function createRecipeStepProduct(
  client: Axios,
  recipeID: string,
  recipeStepID: string,
  input: RecipeStepProductCreationRequestInput,
): Promise<APIResponse<RecipeStepProduct>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<RecipeStepProduct>>(
      `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
