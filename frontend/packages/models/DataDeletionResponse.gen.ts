// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IDataDeletionResponse {
  Successful: boolean;
}

export class DataDeletionResponse implements IDataDeletionResponse {
  Successful: boolean;
  constructor(input: Partial<DataDeletionResponse> = {}) {
    this.Successful = input.Successful || false;
  }
}
