// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidPreparationVessel, 
  APIResponse, 
} from '@dinnerdonebetter/models'; 

export async function archiveValidPreparationVessel(
  client: Axios,
  validPreparationVesselID: string,
	): Promise<  APIResponse <  ValidPreparationVessel >    >   {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete< APIResponse < ValidPreparationVessel  >  >(`/api/v1/valid_preparation_vessels/${validPreparationVesselID}`, {});

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}