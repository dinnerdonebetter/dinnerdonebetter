// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  HouseholdInstrumentOwnership, 
  APIResponse, 
  HouseholdInstrumentOwnershipCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function createHouseholdInstrumentOwnership(
  client: Axios,
  
  input: HouseholdInstrumentOwnershipCreationRequestInput,
): Promise<  APIResponse <  HouseholdInstrumentOwnership >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < HouseholdInstrumentOwnership  >  >(`/api/v1/households/instruments`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}