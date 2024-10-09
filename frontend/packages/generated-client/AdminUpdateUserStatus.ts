// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  UserStatusResponse, 
  APIResponse, 
  UserAccountStatusUpdateInput, 
} from '@dinnerdonebetter/models';

export async function adminUpdateUserStatus(
  client: Axios,
  
  input: UserAccountStatusUpdateInput,
): Promise<  APIResponse <  UserStatusResponse >  >  {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse < UserStatusResponse  >  >(`/api/v1/admin/users/status`, input);
	    if (response.data.error) {
	      reject(new Error(response.data.error.message));
	    }

	    resolve(response.data);
	  });
}