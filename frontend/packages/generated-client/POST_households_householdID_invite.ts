// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  HouseholdInvitation, 
  APIResponse, 
  HouseholdInvitationCreationRequestInput, 
} from '@dinnerdonebetter/models';

export async function pOST_households_householdID_invite(
  client: Axios,
  householdID: string,
  input: HouseholdInvitationCreationRequestInput,
): Promise<  APIResponse <  HouseholdInvitation >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < HouseholdInvitation  >  >(`/api/v1/households/${householdID}/invite`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}