// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipePrepTask, APIResponse, RecipePrepTaskUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateRecipePrepTask(
  client: Axios,
  recipeID: string,
  recipePrepTaskID: string,
  input: RecipePrepTaskUpdateRequestInput,
): Promise<APIResponse<RecipePrepTask>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<RecipePrepTask>>(
      `/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
