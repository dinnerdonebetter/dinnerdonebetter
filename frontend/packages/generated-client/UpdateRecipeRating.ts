// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipeRating, APIResponse, RecipeRatingUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateRecipeRating(
  client: Axios,
  recipeID: string,
  recipeRatingID: string,
  input: RecipeRatingUpdateRequestInput,
): Promise<APIResponse<RecipeRating>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<RecipeRating>>(
      `/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
