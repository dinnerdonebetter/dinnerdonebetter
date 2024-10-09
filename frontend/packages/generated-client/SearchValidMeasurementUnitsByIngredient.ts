// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidMeasurementUnit, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function searchValidMeasurementUnitsByIngredient(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  q: string,
	validIngredientID: string,
	): Promise< QueryFilteredResult< ValidMeasurementUnit >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<ValidMeasurementUnit>  >  >(`/api/v1/valid_measurement_units/by_ingredient/${validIngredientID}`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidMeasurementUnit>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}