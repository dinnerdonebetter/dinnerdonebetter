// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  Household, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function archiveHousehold(
  client: Axios,
  householdID: string,
	): Promise<  APIResponse <  Household >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete< APIResponse < Household  >  >(`/api/v1/households/${householdID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}