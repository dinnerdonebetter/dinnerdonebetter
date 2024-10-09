// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientGroup, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getValidIngredientGroup(
  client: Axios,
  validIngredientGroupID: string,
	): Promise<  APIResponse <  ValidIngredientGroup >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < ValidIngredientGroup  >  >(`/api/v1/valid_ingredient_groups/${validIngredientGroupID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}