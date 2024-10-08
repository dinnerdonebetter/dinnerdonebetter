// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipeStepIngredient, APIResponse, RecipeStepIngredientUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateRecipeStepIngredient(
  client: Axios,
  recipeID: string,
  recipeStepID: string,
  recipeStepIngredientID: string,
  input: RecipeStepIngredientUpdateRequestInput,
): Promise<APIResponse<RecipeStepIngredient>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<RecipeStepIngredient>>(
      `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
