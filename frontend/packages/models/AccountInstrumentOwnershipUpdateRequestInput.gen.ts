// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IAccountInstrumentOwnershipUpdateRequestInput {
  notes: string;
  quantity: number;
  validInstrumentID: string;
}

export class AccountInstrumentOwnershipUpdateRequestInput implements IAccountInstrumentOwnershipUpdateRequestInput {
  notes: string;
  quantity: number;
  validInstrumentID: string;
  constructor(input: Partial<AccountInstrumentOwnershipUpdateRequestInput> = {}) {
    this.notes = input.notes || '';
    this.quantity = input.quantity || 0;
    this.validInstrumentID = input.validInstrumentID || '';
  }
}
