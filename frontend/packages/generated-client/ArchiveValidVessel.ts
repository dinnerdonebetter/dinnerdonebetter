// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidVessel, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function archiveValidVessel(
  client: Axios,
  validVesselID: string,
	): Promise<  APIResponse <  ValidVessel >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete< APIResponse < ValidVessel  >  >(`/api/v1/valid_vessels/${validVesselID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}