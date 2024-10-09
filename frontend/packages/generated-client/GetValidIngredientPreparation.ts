// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientPreparation, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getValidIngredientPreparation(
  client: Axios,
  validIngredientPreparationID: string,
	): Promise<  APIResponse <  ValidIngredientPreparation >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < ValidIngredientPreparation  >  >(`/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}