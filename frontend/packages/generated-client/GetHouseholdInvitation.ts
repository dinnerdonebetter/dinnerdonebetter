// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  HouseholdInvitation, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getHouseholdInvitation(
  client: Axios,
  householdInvitationID: string,
	): Promise<  APIResponse <  HouseholdInvitation >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < HouseholdInvitation  >  >(`/api/v1/household_invitations/${householdInvitationID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}