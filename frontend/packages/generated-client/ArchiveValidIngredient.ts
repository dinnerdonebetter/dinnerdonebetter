// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredient, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function archiveValidIngredient(
  client: Axios,
  validIngredientID: string,
	): Promise<  APIResponse <  ValidIngredient >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete< APIResponse < ValidIngredient  >  >(`/api/v1/valid_ingredients/${validIngredientID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}