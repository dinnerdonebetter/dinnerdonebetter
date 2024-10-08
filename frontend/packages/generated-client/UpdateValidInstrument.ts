// GENERATED CODE, DO NOT EDIT MANUALLY

import { Axios } from 'axios';

import { ValidInstrument, APIResponse, ValidInstrumentUpdateRequestInput } from '@dinnerdonebetter/models';

export async function updateValidInstrument(
  client: Axios,
  validInstrumentID: string,
  input: ValidInstrumentUpdateRequestInput,
): Promise<APIResponse<ValidInstrument>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidInstrument>>(
      `/api/v1/valid_instruments/${validInstrumentID}`,
      input,
    );
    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data);
  });
}
