// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidPreparationInstrument, APIResponse } from '@dinnerdonebetter/models';

export async function getValidPreparationInstrument(
  client: Axios,
  validPreparationVesselID: string,
): Promise<APIResponse<ValidPreparationInstrument>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidPreparationInstrument>>(
      `/api/v1/valid_preparation_instruments/${validPreparationVesselID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
