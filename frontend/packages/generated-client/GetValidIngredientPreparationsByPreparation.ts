// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidIngredientPreparation, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getValidIngredientPreparationsByPreparation(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  validPreparationID: string,
	): Promise< QueryFilteredResult< ValidIngredientPreparation >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<ValidIngredientPreparation>  >  >(`/api/v1/valid_ingredient_preparations/by_preparation/${validPreparationID}`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidIngredientPreparation>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}