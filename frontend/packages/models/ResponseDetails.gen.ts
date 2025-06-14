// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IResponseDetails {
  currentAccountID: string;
  traceID: string;
}

export class ResponseDetails implements IResponseDetails {
  currentAccountID: string;
  traceID: string;
  constructor(input: Partial<ResponseDetails> = {}) {
    this.currentAccountID = input.currentAccountID || '';
    this.traceID = input.traceID || '';
  }
}
