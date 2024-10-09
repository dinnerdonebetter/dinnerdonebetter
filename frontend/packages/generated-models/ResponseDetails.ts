// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IResponseDetails {
  currentHouseholdID: string;
  traceID: string;
}

export class ResponseDetails implements IResponseDetails {
  currentHouseholdID: string;
  traceID: string;
  constructor(input: Partial<ResponseDetails> = {}) {
    this.currentHouseholdID = input.currentHouseholdID = '';
    this.traceID = input.traceID = '';
  }
}
