// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidVessel, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function getRandomValidVessel(
  client: Axios,
  ): Promise<  APIResponse <  ValidVessel >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < ValidVessel  >  >(`/api/v1/valid_vessels/random`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}