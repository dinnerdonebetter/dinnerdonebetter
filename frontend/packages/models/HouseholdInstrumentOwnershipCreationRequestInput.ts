// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IHouseholdInstrumentOwnershipCreationRequestInput {
  validInstrumentID: string;
  belongsToHousehold: string;
  notes: string;
  quantity: number;
}

export class HouseholdInstrumentOwnershipCreationRequestInput
  implements IHouseholdInstrumentOwnershipCreationRequestInput
{
  validInstrumentID: string;
  belongsToHousehold: string;
  notes: string;
  quantity: number;
  constructor(input: Partial<HouseholdInstrumentOwnershipCreationRequestInput> = {}) {
    this.validInstrumentID = input.validInstrumentID || '';
    this.belongsToHousehold = input.belongsToHousehold || '';
    this.notes = input.notes || '';
    this.quantity = input.quantity || 0;
  }
}
