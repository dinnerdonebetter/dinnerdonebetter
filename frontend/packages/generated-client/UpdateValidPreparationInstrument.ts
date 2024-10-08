// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import {
  ValidPreparationInstrument,
  APIResponse,
  ValidPreparationInstrumentUpdateRequestInput,
} from '@dinnerdonebetter/models';

export async function updateValidPreparationInstrument(
  client: Axios,
  validPreparationVesselID: string,
  input: ValidPreparationInstrumentUpdateRequestInput,
): Promise<APIResponse<ValidPreparationInstrument>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidPreparationInstrument>>(
      `/api/v1/valid_preparation_instruments/${validPreparationVesselID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
