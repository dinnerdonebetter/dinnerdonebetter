// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IAccountInstrumentOwnershipCreationRequestInput {
  belongsToAccount: string;
  notes: string;
  quantity: number;
  validInstrumentID: string;
}

export class AccountInstrumentOwnershipCreationRequestInput implements IAccountInstrumentOwnershipCreationRequestInput {
  belongsToAccount: string;
  notes: string;
  quantity: number;
  validInstrumentID: string;
  constructor(input: Partial<AccountInstrumentOwnershipCreationRequestInput> = {}) {
    this.belongsToAccount = input.belongsToAccount || '';
    this.notes = input.notes || '';
    this.quantity = input.quantity || 0;
    this.validInstrumentID = input.validInstrumentID || '';
  }
}
