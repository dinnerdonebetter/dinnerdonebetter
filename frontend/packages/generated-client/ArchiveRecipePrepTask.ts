// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { RecipePrepTask, APIResponse } from '@dinnerdonebetter/models';

export async function archiveRecipePrepTask(
  client: Axios,
  recipeID: string,
  recipePrepTaskID: string,
): Promise<APIResponse<RecipePrepTask>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<RecipePrepTask>>(
      `/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
