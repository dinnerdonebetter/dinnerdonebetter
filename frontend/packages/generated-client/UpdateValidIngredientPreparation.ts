// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientPreparation, 
  APIResponse, 
  ValidIngredientPreparationUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateValidIngredientPreparation(
  client: Axios,
  validIngredientPreparationID: string,
  input: ValidIngredientPreparationUpdateRequestInput,
): Promise<  APIResponse <  ValidIngredientPreparation >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < ValidIngredientPreparation  >  >(`/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}