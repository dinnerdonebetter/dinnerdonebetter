// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  HouseholdInstrumentOwnership, 
  APIResponse, 
  HouseholdInstrumentOwnershipUpdateRequestInput, 
} from '@dinnerdonebetter/models';

export async function updateHouseholdInstrumentOwnership(
  client: Axios,
  householdInstrumentOwnershipID: string,
  input: HouseholdInstrumentOwnershipUpdateRequestInput,
): Promise<  APIResponse <  HouseholdInstrumentOwnership >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse < HouseholdInstrumentOwnership  >  >(`/api/v1/households/instruments/${householdInstrumentOwnershipID}`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}