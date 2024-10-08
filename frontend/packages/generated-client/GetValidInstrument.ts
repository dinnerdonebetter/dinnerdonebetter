// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidInstrument, APIResponse } from '@dinnerdonebetter/models';

export async function getValidInstrument(
  client: Axios,
  validInstrumentID: string,
): Promise<APIResponse<ValidInstrument>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidInstrument>>(
      `/api/v1/valid_instruments/${validInstrumentID}`,
      {},
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
