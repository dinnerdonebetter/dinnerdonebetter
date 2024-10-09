// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  User, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getSelf(
  client: Axios,
  ): Promise<  APIResponse <  User >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < User  >  >(`/api/v1/users/self`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}