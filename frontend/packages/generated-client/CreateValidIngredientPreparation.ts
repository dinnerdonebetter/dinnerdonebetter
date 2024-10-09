// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientPreparation, 
  APIResponse, 
  ValidIngredientPreparationCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createValidIngredientPreparation(
  client: Axios,
  
  input: ValidIngredientPreparationCreationRequestInput,
): Promise<  APIResponse <  ValidIngredientPreparation >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < ValidIngredientPreparation  >  >(`/api/v1/valid_ingredient_preparations`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}