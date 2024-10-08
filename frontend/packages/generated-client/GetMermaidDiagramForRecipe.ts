// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { APIResponse } from '@dinnerdonebetter/models';

export async function getMermaidDiagramForRecipe(client: Axios, recipeID: string): Promise<APIResponse<string>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<string>>(`/api/v1/recipes/${recipeID}/mermaid`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
