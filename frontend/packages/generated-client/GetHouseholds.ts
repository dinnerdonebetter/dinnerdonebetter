// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  Household, 
  QueryFilter,
  QueryFilteredResult,
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getHouseholds(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  ): Promise< QueryFilteredResult< Household >> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Array<Household>  >  >(`/api/v1/households`, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<Household>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}